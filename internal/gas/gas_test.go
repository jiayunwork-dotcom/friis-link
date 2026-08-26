package gas

import "testing"

func TestOxygenRisesToward60(t *testing.T) {
	if err := OxygenRisesToward60(10, 40, StandardAir()); err != nil {
		t.Fatal(err)
	}
}

func TestDryAirLessThanMoist(t *testing.T) {
	if err := DryAirLessThanMoist(22); err != nil {
		t.Fatal(err)
	}
}
