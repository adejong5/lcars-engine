package ha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// LiveClient is a Source backed by a live Home Assistant WebSocket connection.
// It authenticates, seeds from get_states, subscribes to state_changed, and
// reconnects with backoff. Ported from the original ha-websocket.js.
type LiveClient struct {
	host      string
	token     string
	ssl       bool
	log       *slog.Logger
	store     *Store
	connected atomic.Bool
}

// NewLive returns a live client. host may omit the port (defaults to :8123).
func NewLive(host, token string, ssl bool, log *slog.Logger) *LiveClient {
	return &LiveClient{
		host:  normalizeHost(host),
		token: token,
		ssl:   ssl,
		log:   log,
		store: NewStore(),
	}
}

func (c *LiveClient) State(id string) (State, bool) { return c.store.Get(id) }
func (c *LiveClient) All() []State                  { return c.store.All() }

// Connected reports whether the WebSocket is currently authenticated.
func (c *LiveClient) Connected() bool { return c.connected.Load() }

var errAuthInvalid = errors.New("ha: auth_invalid (check HA_TOKEN)")

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
	conn, _, err := websocket.Dial(ctx, c.wsURL(), nil)
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
	c.log.Info("connected to HA", "url", c.wsURL())

	// Seed current states and subscribe to changes.
	if err := wsjson.Write(ctx, conn, cmdMsg{ID: 1, Type: "get_states"}); err != nil {
		return err
	}
	if err := wsjson.Write(ctx, conn, cmdMsg{ID: 2, Type: "subscribe_events", EventType: "state_changed"}); err != nil {
		return err
	}

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
		}
	case "event":
		if m.Event != nil && m.Event.EventType == "state_changed" && m.Event.Data.NewState != nil {
			c.store.Set(*m.Event.Data.NewState)
		}
	}
}

func (c *LiveClient) wsURL() string {
	proto := "ws"
	if c.ssl {
		proto = "wss"
	}
	return fmt.Sprintf("%s://%s/api/websocket", proto, c.host)
}

// normalizeHost appends the default HA port when none is given.
func normalizeHost(host string) string {
	if host == "" || strings.Contains(host, ":") {
		return host
	}
	return host + ":8123"
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

type wsMsg struct {
	Type    string          `json:"type"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Event   *wsEvent        `json:"event"`
	Message string          `json:"message"`
}

type wsEvent struct {
	EventType string `json:"event_type"`
	Data      struct {
		EntityID string `json:"entity_id"`
		NewState *State `json:"new_state"`
	} `json:"data"`
}
