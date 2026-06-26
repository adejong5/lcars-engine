package view

import (
	"math"
	"strconv"
)

// Gauge is a value presented against a range — the model behind PillGauge,
// BarChart, and LevelMeter.
type Gauge struct {
	Label   string
	Display string  // formatted value + unit, or Placeholder when unavailable
	Frac    float64 // position of the value within [Min,Max], clamped to 0..1
	OK      bool    // false when the state is missing/unavailable/non-numeric
}

// GaugeSpec describes how to present one entity as a Gauge.
type GaugeSpec struct {
	Label    string
	Min, Max float64
	Unit     string
	Decimals int
}

// RenderGauge builds a Gauge from a raw state string. Pure.
func RenderGauge(sp GaugeSpec, raw string, known bool) Gauge {
	f, ok := Num(raw)
	if !known || !ok {
		return Gauge{Label: sp.Label, Display: Placeholder, OK: false}
	}
	frac := 0.0
	if sp.Max > sp.Min {
		frac = math.Min(1, math.Max(0, (f-sp.Min)/(sp.Max-sp.Min)))
	}
	disp := strconv.FormatFloat(f, 'f', sp.Decimals, 64) + sp.Unit
	return Gauge{Label: sp.Label, Display: disp, Frac: frac, OK: true}
}
