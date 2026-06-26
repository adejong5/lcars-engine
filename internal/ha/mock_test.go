package ha

import (
	"strconv"
	"testing"
)

func TestMockLazySeedNumeric(t *testing.T) {
	m := NewMock()
	st, ok := m.State("sensor.cpu")
	if !ok {
		t.Fatal("mock should fabricate unknown entity")
	}
	if st.EntityID != "sensor.cpu" {
		t.Fatalf("entity id %q", st.EntityID)
	}
	if _, err := strconv.ParseFloat(st.State, 64); err != nil {
		t.Fatalf("sensor seed %q not numeric", st.State)
	}
	if len(m.All()) != 1 {
		t.Fatalf("seeded entity should appear in All(); got %d", len(m.All()))
	}
	// a second read is stable (not reseeded)
	if st2, _ := m.State("sensor.cpu"); st2.State != st.State {
		t.Fatalf("entity was reseeded: %q != %q", st2.State, st.State)
	}
}

func TestMockBinarySeedIsOnOff(t *testing.T) {
	m := NewMock()
	for i := 0; i < 50; i++ {
		st, _ := m.State("binary_sensor.motion_" + strconv.Itoa(i))
		if st.State != "on" && st.State != "off" {
			t.Fatalf("binary seed %q not on/off", st.State)
		}
	}
}

func TestMockDriftValue(t *testing.T) {
	// numeric stays numeric and within a sane band of the base (+/-5% per step)
	for i := 0; i < 1000; i++ {
		got := mockDriftValue("sensor.x", "100")
		f, err := strconv.ParseFloat(got, 64)
		if err != nil {
			t.Fatalf("drift produced non-numeric %q", got)
		}
		if f < 90 || f > 110 {
			t.Fatalf("drift %v outside a plausible band of 100", f)
		}
	}
	// binary stays on/off
	for i := 0; i < 1000; i++ {
		if got := mockDriftValue("binary_sensor.x", "on"); got != "on" && got != "off" {
			t.Fatalf("binary drift %q not on/off", got)
		}
	}
	// non-numeric, non-binary passes through unchanged
	if got := mockDriftValue("sensor.text", "idle"); got != "idle" {
		t.Fatalf("non-numeric drift changed value: %q", got)
	}
}
