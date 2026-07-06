package view

import "strconv"

// Readout is one horizontal instrument line: a left cap, the entity label on
// the bar colour, the state as large text on a black cutout (the text-bar
// treatment), a colour-coded fill bar, and a pill end cap carrying the unit.
type Readout struct {
	Label string
	Value string  // formatted number (the unit lives in the end cap)
	Unit  string
	Frac  float64 // value position in [Min,Max], clamped 0..1
	Color string  // coded colour for the value + fill (nearest stop)
	OK    bool
}

// ReadoutSpec describes how to present one entity as a Readout.
type ReadoutSpec struct {
	Label    string
	Min, Max float64
	Unit     string
	Decimals int
	Stops    []string // coded colours low→high (like a meter's gradient stops)
}

// RenderReadout builds a Readout from a raw state string. Pure.
func RenderReadout(sp ReadoutSpec, raw string, known bool) Readout {
	ro := Readout{Label: sp.Label, Unit: sp.Unit, Value: Placeholder}
	if len(sp.Stops) > 0 {
		ro.Color = sp.Stops[0]
	}
	f, ok := Num(raw)
	if !known || !ok {
		return ro
	}
	frac := 0.0
	if sp.Max > sp.Min {
		frac = (f - sp.Min) / (sp.Max - sp.Min)
	}
	if frac < 0 {
		frac = 0
	} else if frac > 1 {
		frac = 1
	}
	ro.Value = strconv.FormatFloat(f, 'f', sp.Decimals, 64)
	ro.Frac = frac
	ro.OK = true
	if n := len(sp.Stops); n > 0 {
		ro.Color = sp.Stops[int(frac*float64(n-1)+0.5)]
	}
	return ro
}
