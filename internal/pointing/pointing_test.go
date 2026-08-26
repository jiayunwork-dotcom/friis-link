package pointing

import "testing"

func TestPointingLossFromBeamwidth(t *testing.T) {
	if err := NarrowerBeamHurtsMore(10, 25, 1.5); err != nil {
		t.Fatal(err)
	}
	if err := OrthogonalRejected(); err != nil {
		t.Fatal(err)
	}
	loss, err := PolarizationMismatch(45)
	if err != nil {
		t.Fatal(err)
	}
	if loss < 2.5 || loss > 3.5 {
		t.Fatalf("45° mismatch should be ~3 dB, got %g", loss)
	}
}
