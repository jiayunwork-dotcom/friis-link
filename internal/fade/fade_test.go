package fade

import (
	"math"
	"testing"
)

func TestRayleighOutageFallsWithMargin(t *testing.T) {
	o1, err := RayleighOutage(0)
	if err != nil {
		t.Fatal(err)
	}
	o2, err := RayleighOutage(10)
	if err != nil {
		t.Fatal(err)
	}
	if !(o2 < o1) {
		t.Fatalf("more margin should cut outage, %g vs %g", o2, o1)
	}
	req, err := RequiredRayleigh(0.01)
	if err != nil {
		t.Fatal(err)
	}
	back, err := RayleighOutage(req)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(back-0.01) > 1e-6 {
		t.Fatalf("round-trip outage %g", back)
	}
	c, err := CombineMargins(6, 8)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(c-10) > 1e-12 {
		t.Fatalf("rss %g", c)
	}
	if err := RequiredMatchesOutage(0.01); err != nil {
		t.Fatal(err)
	}
	if err := LargerMarginLowerOutage(3, 12); err != nil {
		t.Fatal(err)
	}
}
