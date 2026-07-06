// Package view turns raw Home Assistant state strings into display-ready values
// for the templates. Everything here is pure (no I/O), so it is fully
// unit-tested and keeps the eventual template layer to plain markup binding.
package view

import (
	"math"
	"strconv"
	"strings"
)

// Placeholder is shown when a value is missing or unavailable (matches the
// original dashboards' "----").
const Placeholder = "----"

// AlertClass is the LCARS styling applied by Hot/Low when a threshold trips.
const AlertClass = "font-red blink"

// Formatter turns a raw state string into display text.
type Formatter func(raw string) string

// Classifier returns extra CSS classes for a raw value (e.g. alert styling).
type Classifier func(raw string) string

// Num parses a raw state into a float. ok is false for empty/unknown/
// unavailable values or anything non-numeric.
func Num(raw string) (float64, bool) {
	s := strings.TrimSpace(raw)
	switch s {
	case "", "unknown", "unavailable":
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, false
	}
	return f, true
}

// Round formats a numeric state as a whole number (the old pct helper);
// non-numeric input is returned unchanged.
func Round(raw string) string {
	if f, ok := Num(raw); ok {
		return strconv.FormatFloat(math.Round(f), 'f', -1, 64)
	}
	return raw
}

// Fixed returns a Formatter that renders a numeric state to n decimal places
// (the old one == Fixed(1)); non-numeric input is returned unchanged.
func Fixed(n int) Formatter {
	return func(raw string) string {
		if f, ok := Num(raw); ok {
			return strconv.FormatFloat(f, 'f', n, 64)
		}
		return raw
	}
}

// Map returns a Formatter that substitutes known raw states with display text
// (e.g. binary_sensor on/off -> OPEN/CLOSED); unknown values pass through.
func Map(subs map[string]string) Formatter {
	return func(raw string) string {
		if v, ok := subs[strings.TrimSpace(raw)]; ok {
			return v
		}
		return raw
	}
}

// Map2 returns a Classifier that maps a raw state to a CSS class (e.g. an
// "on" door -> AlertClass); unmatched states get no class.
func Map2(subs map[string]string) Classifier {
	return func(raw string) string { return subs[strings.TrimSpace(raw)] }
}

// Hot returns a Classifier that flags AlertClass when the value is >= limit.
func Hot(limit float64) Classifier {
	return func(raw string) string {
		if f, ok := Num(raw); ok && f >= limit {
			return AlertClass
		}
		return ""
	}
}

// Low returns a Classifier that flags AlertClass when the value is < limit.
func Low(limit float64) Classifier {
	return func(raw string) string {
		if f, ok := Num(raw); ok && f < limit {
			return AlertClass
		}
		return ""
	}
}

// unavailable reports whether a raw state means "no usable value".
func unavailable(raw string) bool {
	switch strings.TrimSpace(raw) {
	case "", "unknown", "unavailable":
		return true
	}
	return false
}
