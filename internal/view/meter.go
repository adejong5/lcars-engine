package view

import "strconv"

// Meter is the precomputed geometry for a LevelMeter: a vertical stack of
// pill segments that light from the bottom up, each showing its slice of a
// gradient. Ported from the SVG LevelMeter component; geometry is computed here
// so the template just draws it.
type Meter struct {
	GID    string // unique gradient id (multiple meters per page)
	VW, VH float64
	RectW  float64
	Stops  []MeterStop
	Segs   []MeterSeg
}

// MeterStop is one gradient colour stop (bottom→top).
type MeterStop struct {
	Offset float64 // percent
	Color  string  // CSS colour
}

// MeterSeg is one pill segment. Lit segments carry their own solid colour (the
// nearest gradient stop for their position): an SVG gradient would need a
// url(#id) reference, which <base href> re-roots and breaks on every page
// that isn't the base.
type MeterSeg struct {
	Y, H, R float64
	Lit     bool
	Color   string
}

// MeterSpec describes how to present one entity as a labelled segmented meter.
type MeterSpec struct {
	Label    string
	Min, Max float64
	Unit     string
	Decimals int
	Segments int
	Stops    []string // gradient colours bottom→top
}

// MeterBlock is a Meter plus its readout — the model for the meter-block
// fragment (segmented gauge + value + label).
type MeterBlock struct {
	Meter   Meter
	Label   string
	Display string
	OK      bool
}

// RenderMeterBlock builds a MeterBlock from a raw state string. id must be
// unique per page (it namespaces the SVG gradient). Pure.
func RenderMeterBlock(id string, sp MeterSpec, raw string, known bool) MeterBlock {
	f, ok := Num(raw)
	if !known || !ok {
		return MeterBlock{Meter: RenderMeter(id, sp.Min, sp.Min, sp.Max, sp.Segments, 3, sp.Stops), Label: sp.Label, Display: Placeholder}
	}
	disp := strconv.FormatFloat(f, 'f', sp.Decimals, 64) + sp.Unit
	return MeterBlock{Meter: RenderMeter(id, f, sp.Min, sp.Max, sp.Segments, 3, sp.Stops), Label: sp.Label, Display: disp, OK: true}
}

// RenderMeter computes a Meter for value within [min,max], with the given
// segment count, inter-segment gap, and gradient stops (CSS colours, bottom→top).
func RenderMeter(id string, value, min, max float64, segments, gap int, stops []string) Meter {
	const vw, vh = 40.0, 280.0

	frac := 0.0
	if max > min {
		frac = (value - min) / (max - min)
	}
	if frac < 0 {
		frac = 0
	} else if frac > 1 {
		frac = 1
	}
	lit := int(frac*float64(segments) + 0.5) // round to nearest segment

	segH := (vh - float64(segments-1)*float64(gap)) / float64(segments)
	r := segH / 2
	segs := make([]MeterSeg, segments)
	for i := 0; i < segments; i++ {
		fromBottom := segments - 1 - i
		color := ""
		if len(stops) > 0 {
			pos := 0.0
			if segments > 1 {
				pos = float64(fromBottom) / float64(segments-1)
			}
			color = stops[int(pos*float64(len(stops)-1)+0.5)]
		}
		segs[i] = MeterSeg{Y: float64(i) * (segH + float64(gap)), H: segH, R: r, Lit: fromBottom < lit, Color: color}
	}

	ms := make([]MeterStop, len(stops))
	for i, c := range stops {
		off := 0.0
		if len(stops) > 1 {
			off = float64(i) / float64(len(stops)-1) * 100
		}
		ms[i] = MeterStop{Offset: off, Color: c}
	}

	return Meter{GID: id, VW: vw, VH: vh, RectW: vw - 4, Stops: ms, Segs: segs}
}
