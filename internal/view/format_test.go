package view

import "testing"

func TestNum(t *testing.T) {
	ok := []struct {
		in   string
		want float64
	}{
		{"42", 42}, {"72.68", 72.68}, {"-3.5", -3.5}, {" 10 ", 10}, {"0", 0},
	}
	for _, c := range ok {
		if got, k := Num(c.in); !k || got != c.want {
			t.Fatalf("Num(%q) = %v,%v; want %v,true", c.in, got, k, c.want)
		}
	}
	for _, in := range []string{"", "unknown", "unavailable", "on", "abc", " "} {
		if _, k := Num(in); k {
			t.Fatalf("Num(%q) should be !ok", in)
		}
	}
}

func TestRound(t *testing.T) {
	cases := map[string]string{"72.68": "73", "12.0": "12", "4.4": "4", "-1.6": "-2"}
	for in, want := range cases {
		if got := Round(in); got != want {
			t.Fatalf("Round(%q) = %q, want %q", in, got, want)
		}
	}
	if got := Round("unavailable"); got != "unavailable" {
		t.Fatalf("Round passthrough for non-numeric: %q", got)
	}
}

func TestFixed(t *testing.T) {
	f := Fixed(1)
	if got := f("72.68"); got != "72.7" {
		t.Fatalf("Fixed(1)(72.68) = %q, want 72.7", got)
	}
	if got := f("5"); got != "5.0" {
		t.Fatalf("Fixed(1)(5) = %q, want 5.0", got)
	}
	if got := Fixed(2)("abc"); got != "abc" {
		t.Fatalf("Fixed passthrough for non-numeric: %q", got)
	}
}

func TestHotLow(t *testing.T) {
	hot := Hot(85)
	if hot("90") != AlertClass {
		t.Fatal("Hot(85)(90) should alert")
	}
	if hot("85") != AlertClass {
		t.Fatal("Hot is inclusive at the limit")
	}
	if hot("84") != "" {
		t.Fatal("Hot(85)(84) should not alert")
	}
	if hot("unavailable") != "" {
		t.Fatal("Hot on unavailable should not alert")
	}

	low := Low(100)
	if low("99") != AlertClass {
		t.Fatal("Low(100)(99) should alert")
	}
	if low("100") != "" {
		t.Fatal("Low is exclusive at the limit")
	}
}
