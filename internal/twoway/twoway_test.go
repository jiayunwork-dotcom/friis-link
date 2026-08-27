package twoway

import (
	"math"
	"testing"

	"friis-link/internal/propagation"
)

func TestCrossoverAndFarFieldPenalty(t *testing.T) {
	lam, err := propagation.Wavelength(2.45e9)
	if err != nil {
		t.Fatal(err)
	}
	dc, err := Crossover(20, 20, lam)
	if err != nil {
		t.Fatal(err)
	}
	if dc <= 0 {
		t.Fatal("crossover")
	}
	far, err := InTwoRay(3*dc, 20, 20, 2.45e9)
	if err != nil {
		t.Fatal(err)
	}
	if !far {
		t.Fatal("3 dc should be two-ray")
	}
	near, err := InTwoRay(dc/3, 20, 20, 2.45e9)
	if err != nil {
		t.Fatal(err)
	}
	if near {
		t.Fatal("dc/3 should still follow Friis")
	}
	if DoubleDistancePenaltyDB(true) != 12 || DoubleDistancePenaltyDB(false) != 6 {
		t.Fatal("penalties")
	}
	pe1, err := PlaneEarthLossDB(10000, 20, 20)
	if err != nil {
		t.Fatal(err)
	}
	pe2, err := PlaneEarthLossDB(20000, 20, 20)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs((pe2-pe1)-40*math.Log10(2)) > 1e-9 {
		t.Fatalf("two-ray d×2 should add 40 log10(2) dB, got %g", pe2-pe1)
	}
}

func TestTwoRayBeyondCrossoverDrops12(t *testing.T) {
	extra, err := FarFieldPrDrop(40000, 20, 20, 2.45e9)
	if err != nil {
		t.Fatal(err)
	}
	if extra < 5.5 {
		t.Fatalf("two-ray should cost ~6 dB extra vs Friis on d×2, got %g", extra)
	}
	if err := ExcessGrowsPastCrossover(20, 20, 2.45e9); err != nil {
		t.Fatal(err)
	}
}
