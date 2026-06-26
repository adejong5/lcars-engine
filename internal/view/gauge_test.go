package view

import "testing"

func TestRenderGaugeFrac(t *testing.T) {
	sp := GaugeSpec{Label: "Temp", Min: 50, Max: 90, Unit: "°F", Decimals: 1}

	mid := RenderGauge(sp, "70", true)
	if !mid.OK || mid.Frac != 0.5 {
		t.Fatalf("mid frac = %v ok=%v, want 0.5", mid.Frac, mid.OK)
	}
	if mid.Display != "70.0°F" {
		t.Fatalf("display = %q, want 70.0°F", mid.Display)
	}

	// below min and above max clamp to 0 and 1
	if g := RenderGauge(sp, "40", true); g.Frac != 0 {
		t.Fatalf("below-min frac = %v, want 0", g.Frac)
	}
	if g := RenderGauge(sp, "120", true); g.Frac != 1 {
		t.Fatalf("above-max frac = %v, want 1", g.Frac)
	}
}

func TestRenderGaugeUnavailable(t *testing.T) {
	sp := GaugeSpec{Label: "Temp", Min: 0, Max: 100}
	for _, raw := range []string{"", "unknown", "unavailable", "abc"} {
		g := RenderGauge(sp, raw, true)
		if g.OK || g.Display != Placeholder {
			t.Fatalf("RenderGauge(%q) = %+v, want placeholder", raw, g)
		}
	}
	if g := RenderGauge(sp, "50", false); g.OK {
		t.Fatalf("unknown entity should not be OK: %+v", g)
	}
}

func TestRenderGaugeDegenerateRange(t *testing.T) {
	// Max <= Min must not divide by zero; frac stays 0
	g := RenderGauge(GaugeSpec{Label: "X", Min: 10, Max: 10}, "10", true)
	if !g.OK || g.Frac != 0 {
		t.Fatalf("degenerate range = %+v, want OK frac 0", g)
	}
}
