package ha

import (
	"encoding/json"
	"io"
	"log/slog"
	"testing"
)

func testClient() *LiveClient {
	return &LiveClient{store: NewStore(), log: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func TestHandleGetStatesResult(t *testing.T) {
	c := testClient()
	raw := `{"id":1,"type":"result","success":true,"result":[
	  {"entity_id":"switch.kitchen_spare","state":"off","attributes":{"friendly_name":"Kitchen Spare"}},
	  {"entity_id":"sensor.kitchen_current_temperature","state":"72.68","attributes":{}}
	]}`
	var m wsMsg
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	c.handle(m)

	if st, ok := c.store.Get("switch.kitchen_spare"); !ok || st.State != "off" {
		t.Fatalf("kitchen switch not seeded: %+v ok=%v", st, ok)
	}
	if st, _ := c.store.Get("sensor.kitchen_current_temperature"); st.State != "72.68" {
		t.Fatalf("temperature not seeded: %q", st.State)
	}
}

func TestHandleStateChangedEvent(t *testing.T) {
	c := testClient()
	c.store.Set(State{EntityID: "switch.kitchen_spare", State: "off"})

	raw := `{"id":2,"type":"event","event":{"event_type":"state_changed","data":{
	  "entity_id":"switch.kitchen_spare",
	  "new_state":{"entity_id":"switch.kitchen_spare","state":"on","attributes":{}},
	  "old_state":{"entity_id":"switch.kitchen_spare","state":"off","attributes":{}}
	}}}`
	var m wsMsg
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	c.handle(m)

	if st, _ := c.store.Get("switch.kitchen_spare"); st.State != "on" {
		t.Fatalf("event did not update state, got %q", st.State)
	}
}

func TestHandleIgnoresOtherResults(t *testing.T) {
	c := testClient()
	// the subscribe_events ack (id 2) must not be treated as states
	raw := `{"id":2,"type":"result","success":true,"result":null}`
	var m wsMsg
	_ = json.Unmarshal([]byte(raw), &m)
	c.handle(m)
	if len(c.store.All()) != 0 {
		t.Fatalf("non-get_states result populated the store: %v", c.store.All())
	}
}

func TestBuildWSURL(t *testing.T) {
	if got := BuildWSURL("192.168.2.100", false); got != "ws://192.168.2.100:8123/api/websocket" {
		t.Fatalf("standalone url: %q", got)
	}
	if got := BuildWSURL("ha.local:8123", true); got != "wss://ha.local:8123/api/websocket" {
		t.Fatalf("ssl url: %q", got)
	}
	if SupervisorWSURL != "ws://supervisor/core/websocket" {
		t.Fatalf("supervisor url: %q", SupervisorWSURL)
	}
}

func TestNormalizeHost(t *testing.T) {
	if got := normalizeHost("192.168.2.100"); got != "192.168.2.100:8123" {
		t.Fatalf("default port not applied: %q", got)
	}
	if got := normalizeHost("ha.local:8123"); got != "ha.local:8123" {
		t.Fatalf("explicit port changed: %q", got)
	}
	if got := normalizeHost(""); got != "" {
		t.Fatalf("empty host should stay empty: %q", got)
	}
}

func TestBackoff(t *testing.T) {
	cases := map[int]string{0: "3s", 1: "6s", 2: "12s", 3: "24s", 4: "48s", 5: "1m0s", 99: "1m0s"}
	for attempt, want := range cases {
		if got := backoff(attempt).String(); got != want {
			t.Fatalf("backoff(%d) = %s, want %s", attempt, got, want)
		}
	}
}
