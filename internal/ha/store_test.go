package ha

import "testing"

func TestStoreSetGetHas(t *testing.T) {
	s := NewStore()
	if _, ok := s.Get("sensor.x"); ok {
		t.Fatal("empty store should not have sensor.x")
	}
	if s.Has("sensor.x") {
		t.Fatal("empty store Has should be false")
	}

	s.Set(State{EntityID: "sensor.x", State: "42"})
	st, ok := s.Get("sensor.x")
	if !ok || st.State != "42" {
		t.Fatalf("got %+v ok=%v, want state 42", st, ok)
	}
	if !s.Has("sensor.x") {
		t.Fatal("Has should be true after Set")
	}

	s.Set(State{EntityID: "sensor.x", State: "43"}) // overwrite
	if st, _ := s.Get("sensor.x"); st.State != "43" {
		t.Fatalf("overwrite failed: %s", st.State)
	}
}

func TestStoreAllSorted(t *testing.T) {
	s := NewStore()
	s.Set(State{EntityID: "sensor.b"})
	s.Set(State{EntityID: "sensor.a"})
	s.Set(State{EntityID: "sensor.c"})

	all := s.All()
	if len(all) != 3 {
		t.Fatalf("want 3, got %d", len(all))
	}
	if all[0].EntityID != "sensor.a" || all[1].EntityID != "sensor.b" || all[2].EntityID != "sensor.c" {
		t.Fatalf("not sorted by entity_id: %v", all)
	}
}
