package view

import (
	"math"
	"strconv"
)

// Bar is one horizontal bar in a BarChart: a value as a fraction of a max, with
// a formatted readout and a fill colour.
type Bar struct {
	Label   string
	Display string  // formatted value + unit, or Placeholder
	Frac    float64 // value / max, clamped to 0..1
	Color   string  // CSS colour for the fill
}

// BarSpec describes how to present one entity as a Bar.
type BarSpec struct {
	Label    string
	Max      float64
	Unit     string
	Decimals int
	Color    string
}

// RenderBar builds a Bar from a raw state string. Pure.
func RenderBar(sp BarSpec, raw string, known bool) Bar {
	f, ok := Num(raw)
	if !known || !ok {
		return Bar{Label: sp.Label, Display: Placeholder, Color: sp.Color}
	}
	frac := 0.0
	if sp.Max > 0 {
		frac = math.Min(1, math.Max(0, f/sp.Max))
	}
	return Bar{
		Label:   sp.Label,
		Display: strconv.FormatFloat(f, 'f', sp.Decimals, 64) + sp.Unit,
		Frac:    frac,
		Color:   sp.Color,
	}
}
