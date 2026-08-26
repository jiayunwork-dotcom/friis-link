package rain

import (
	"math"
	"testing"
)

func TestHigherBandMoreRain(t *testing.T) {
	if err := HigherBandMoreRain(25, 10); err != nil {
		t.Fatal(err)
	}
	if err := ReductionShrinksPath(0.01, 1.2, 25, 10); err != nil {
		t.Fatal(err)
	}
}

func TestSpecificScalesWithRatePowerAndZeroIsZero(t *testing.T) {
	z, err := Specific(0.01, 1.1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if z != 0 {
		t.Fatal("zero rain")
	}
	a, b, err := DoubleRate(0.01, 1.0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(b/a-2) > 1e-12 {
		t.Fatalf("alpha=1 should double with rate, %g", b/a)
	}
	p, err := Path(0.01, 1.0, 10, 5)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(p-0.5) > 1e-12 {
		t.Fatalf("path %g", p)
	}
	pe, err := PathEffective(0.01, 1.0, 10, 5, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(pe-0.25) > 1e-12 {
		t.Fatalf("effective path %g", pe)
	}
}
