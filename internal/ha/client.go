package ha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// LiveClient is a Source backed by a live Home Assistant WebSocket connection.
// It authenticates, seeds from get_states, subscribes to state_changed, and
// reconnects with backoff. Ported from the original ha-websocket.js.
type LiveClient struct {
	url       string
	token     string
	log       *slog.Logger
	store     *Store
	connected atomic.Bool

	mu     sync.Mutex      // serializes writes and guards conn
	conn   *websocket.Conn // active connection; nil when disconnected
	nextID atomic.Int64    // command id counter (ids 1,2 reserved for setup)

	pendingMu sync.Mutex             // guards pending
	pending   map[int]chan callResult // command id → waiter for HA's result
}

// callResult carries HA's reply to a command back to the waiting caller.
type callResult struct {
	success bool
	err     string
}

// NewLive returns a live client for the given WebSocket URL. Build the URL with
// BuildWSURL (standalone) or use SupervisorWSURL (add-on).
func NewLive(url, token string, log *slog.Logger) *LiveClient {
	c := &LiveClient{
		url:     url,
		token:   token,
		log:     log,
		store:   NewStore(),
		pending: make(map[int]chan callResult),
	}
	c.nextID.Store(2) // 1 = get_states, 2 = subscribe_events
	return c
}

func (c *LiveClient) State(id string) (State, bool) { return c.store.Get(id) }
func (c *LiveClient) All() []State                  { return c.store.All() }

// Connected reports whether the WebSocket is currently authenticated.
func (c *LiveClient) Connected() bool { return c.connected.Load() }

// ErrNotConnected is returned by CallService when there is no live connection.
var ErrNotConnected = errors.New("ha: not connected to Home Assistant")

var errAuthInvalid = errors.New("ha: auth_invalid (check HA_TOKEN)")

// CallService invokes a Home Assistant service and waits for HA's result,
// returning an error if HA reports failure, the connection is down, or the
// reply times out. The resulting state still arrives via the state_changed
// subscription; this just confirms the command was accepted.
func (c *LiveClient) CallService(domain, service string, data, target map[string]any) error {
	id := int(c.nextID.Add(1))
	ch := make(chan callResult, 1)
	c.pendingMu.Lock()
	c.pending[id] = ch
	c.pendingMu.Unlock()
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	msg := callServiceMsg{
		ID:          id,
		Type:        "call_service",
		Domain:      domain,
		Service:     service,
		ServiceData: data,
		Target:      target,
	}

	c.mu.Lock()
	if c.conn == nil {
		c.mu.Unlock()
		return ErrNotConnected
	}
	wctx, wcancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := wsjson.Write(wctx, c.conn, msg)
	wcancel()
	c.mu.Unlock()
	if err != nil {
		return err
	}

	select {
	case res := <-ch:
		if !res.success {
			return fmt.Errorf("ha: %s.%s failed: %s", domain, service, res.err)
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("ha: timed out awaiting %s.%s result", domain, service)
	}
}

// Run maintains the connection until ctx is cancelled, reconnecting with
// exponential backoff (3s..60s). A bad token stops the loop — retrying it is
// pointless and noisy.
func (c *LiveClient) Run(ctx context.Context) {
	attempt := 0
	for ctx.Err() == nil {
		err := c.connectOnce(ctx)
		if ctx.Err() != nil {
			return
		}
		if errors.Is(err, errAuthInvalid) {
			c.log.Error("HA auth invalid; not reconnecting", "err", err)
			return
		}
		delay := backoff(attempt)
		attempt++
		c.log.Warn("HA disconnected; reconnecting", "in", delay.String(), "err", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}
	}
}

func (c *LiveClient) connectOnce(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, c.url, nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "")
	conn.SetReadLimit(32 << 20) // get_states can be large (hundreds of entities)

	// Auth handshake: auth_required -> auth -> auth_ok | auth_invalid.
	var m wsMsg
	if err := wsjson.Read(ctx, conn, &m); err != nil {
		return err
	}
	if err := wsjson.Write(ctx, conn, authMsg{Type: "auth", AccessToken: c.token}); err != nil {
		return err
	}
	if err := wsjson.Read(ctx, conn, &m); err != nil {
		return err
	}
	switch m.Type {
	case "auth_ok":
		// authenticated
	case "auth_invalid":
		return errAuthInvalid
	default:
		return fmt.Errorf("ha: unexpected auth reply %q", m.Type)
	}

	c.connected.Store(true)
	defer c.connected.Store(false)
	c.log.Info("connected to HA", "url", c.url)

	// Seed current states and subscribe to changes.
	if err := wsjson.Write(ctx, conn, cmdMsg{ID: 1, Type: "get_states"}); err != nil {
		return err
	}
	if err := wsjson.Write(ctx, conn, cmdMsg{ID: 2, Type: "subscribe_events", EventType: "state_changed"}); err != nil {
		return err
	}

	// Publish the connection so CallService can write to it.
	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		c.conn = nil
		c.mu.Unlock()
		c.failPending() // unblock any in-flight CallService waiters
	}()

	for {
		var msg wsMsg
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			return err
		}
		c.handle(msg)
	}
}

// handle applies one decoded message to the store. It performs no I/O, so it is
// unit-tested directly with captured fixtures.
func (c *LiveClient) handle(m wsMsg) {
	switch m.Type {
	case "result":
		if m.ID == 1 && len(m.Result) > 0 { // the get_states reply
			var states []State
			if err := json.Unmarshal(m.Result, &states); err != nil {
				c.log.Warn("get_states decode failed", "err", err)
				return
			}
			for _, st := range states {
				c.store.Set(st)
			}
			c.log.Info("seeded states from HA", "count", len(states))
			return
		}
		// Deliver to a waiting CallService, if any.
		c.pendingMu.Lock()
		ch, waiting := c.pending[m.ID]
		c.pendingMu.Unlock()
		if waiting {
			ch <- callResult{success: m.Success != nil && *m.Success, err: errText(m.Error)}
			return
		}
		if m.Success != nil && !*m.Success {
			c.log.Warn("HA command failed", "id", m.ID, "error", errText(m.Error))
		}
	case "event":
		if m.Event != nil && m.Event.EventType == "state_changed" && m.Event.Data.NewState != nil {
			c.store.Set(*m.Event.Data.NewState)
		}
	}
}

// SupervisorWSURL is the HA WebSocket endpoint when running as an add-on: the
// Supervisor proxies to Home Assistant Core, authenticated with SUPERVISOR_TOKEN.
const SupervisorWSURL = "ws://supervisor/core/websocket"

// BuildWSURL builds the standalone HA WebSocket URL. host may omit the port
// (defaults to :8123).
func BuildWSURL(host string, ssl bool) string {
	proto := "ws"
	if ssl {
		proto = "wss"
	}
	return fmt.Sprintf("%s://%s/api/websocket", proto, normalizeHost(host))
}

// normalizeHost appends the default HA port when none is given.
func normalizeHost(host string) string {
	if host == "" || strings.Contains(host, ":") {
		return host
	}
	return host + ":8123"
}

// failPending notifies any in-flight CallService waiters that the connection
// dropped, so they return promptly instead of waiting out the timeout.
func (c *LiveClient) failPending() {
	c.pendingMu.Lock()
	for id, ch := range c.pending {
		select {
		case ch <- callResult{err: "connection lost"}:
		default:
		}
		delete(c.pending, id)
	}
	c.pendingMu.Unlock()
}

// backoff returns 3s, 6s, 12s, … capped at 60s.
func backoff(attempt int) time.Duration {
	d := (3 * time.Second) << attempt
	if d <= 0 || d > 60*time.Second {
		return 60 * time.Second
	}
	return d
}

// ── wire types ──────────────────────────────────────────────

type authMsg struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type cmdMsg struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type,omitempty"`
}

type callServiceMsg struct {
	ID          int            `json:"id"`
	Type        string         `json:"type"`
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      map[string]any `json:"target,omitempty"`
}

type wsMsg struct {
	Type    string          `json:"type"`
	ID      int             `json:"id"`
	Success *bool           `json:"success"`
	Result  json.RawMessage `json:"result"`
	Event   *wsEvent        `json:"event"`
	Error   *wsError        `json:"error"`
	Message string          `json:"message"`
}

type wsError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errText(e *wsError) string {
	if e == nil {
		return ""
	}
	return e.Code + ": " + e.Message
}

type wsEvent struct {
	EventType string `json:"event_type"`
	Data      struct {
		EntityID string `json:"entity_id"`
		NewState *State `json:"new_state"`
	} `json:"data"`
}
