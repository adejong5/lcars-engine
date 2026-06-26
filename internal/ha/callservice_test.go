package ha

import "testing"

func TestMockCallServiceOnOffToggle(t *testing.T) {
	m := NewMock()
	tgt := map[string]any{"entity_id": "switch.test"}

	if err := m.CallService("switch", "turn_on", nil, tgt); err != nil {
		t.Fatal(err)
	}
	if st, _ := m.State("switch.test"); st.State != "on" {
		t.Fatalf("turn_on => %q, want on", st.State)
	}

	if err := m.CallService("switch", "turn_off", nil, tgt); err != nil {
		t.Fatal(err)
	}
	if st, _ := m.State("switch.test"); st.State != "off" {
		t.Fatalf("turn_off => %q, want off", st.State)
	}

	_ = m.CallService("switch", "toggle", nil, tgt)
	if st, _ := m.State("switch.test"); st.State != "on" {
		t.Fatalf("toggle => %q, want on", st.State)
	}
}

func TestTargetEntities(t *testing.T) {
	if got := targetEntities(map[string]any{"entity_id": "switch.a"}); len(got) != 1 || got[0] != "switch.a" {
		t.Fatalf("string target: %v", got)
	}
	if got := targetEntities(map[string]any{"entity_id": []any{"switch.a", "switch.b"}}); len(got) != 2 {
		t.Fatalf("list target: %v", got)
	}
	if got := targetEntities(nil); got != nil {
		t.Fatalf("nil target: %v", got)
	}
	if got := targetEntities(map[string]any{}); got != nil {
		t.Fatalf("missing entity_id: %v", got)
	}
}
