package view

import "testing"

func TestRenderReadout(t *testing.T) {
	sp := ReadoutSpec{Label: "Deck 2", Min: 50, Max: 90, Unit: "°F", Decimals: 0,
		Stops: []string{"var(--ice)", "var(--gold)", "var(--mars)"}}

	mid := RenderReadout(sp, "70", true)
	if !mid.OK || mid.Frac != 0.5 {
		t.Fatalf("mid frac = %v ok=%v, want 0.5", mid.Frac, mid.OK)
	}
	if mid.Value != "70" || mid.Unit != "°F" {
		t.Fatalf("value/unit = %q %q, want 70 °F", mid.Value, mid.Unit)
	}
	if mid.Color != "var(--gold)" {
		t.Fatalf("mid colour = %q, want the middle stop", mid.Color)
	}

	// clamped ends pick the outer stops
	if ro := RenderReadout(sp, "40", true); ro.Frac != 0 || ro.Color != "var(--ice)" {
		t.Fatalf("below-min = %+v, want frac 0 / first stop", ro)
	}
	if ro := RenderReadout(sp, "120", true); ro.Frac != 1 || ro.Color != "var(--mars)" {
		t.Fatalf("above-max = %+v, want frac 1 / last stop", ro)
	}
}

func TestRenderReadoutUnavailable(t *testing.T) {
	sp := ReadoutSpec{Label: "X", Min: 0, Max: 100, Unit: "%",
		Stops: []string{"var(--ice)", "var(--mars)"}}
	for _, raw := range []string{"", "unknown", "unavailable", "abc"} {
		ro := RenderReadout(sp, raw, true)
		if ro.OK || ro.Value != Placeholder {
			t.Fatalf("RenderReadout(%q) = %+v, want placeholder", raw, ro)
		}
		if ro.Color != "var(--ice)" {
			t.Fatalf("unavailable colour = %q, want the first stop", ro.Color)
		}
	}
	if ro := RenderReadout(sp, "50", false); ro.OK {
		t.Fatalf("unknown entity should not be OK: %+v", ro)
	}
}
