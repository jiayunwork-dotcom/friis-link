package model

import (
	"math"
	"testing"
)

func TestDBToLinearRoundTrip(t *testing.T) {
	cases := []struct {
		db     float64
		linear float64
	}{
		{0, 1},
		{10, 10},
		{30, 1000},
		{-10, 0.1},
	}
	for _, c := range cases {
		got := FromDB(c.db)
		if math.Abs(got-c.linear) > 1e-12 {
			t.Errorf("FromDB(%g) = %g, want %g", c.db, got, c.linear)
		}
		back, err := dB(c.linear)
		if err != nil {
			t.Errorf("dB(%g) returned error: %v", c.linear, err)
			continue
		}
		if math.Abs(back-c.db) > 1e-9 {
			t.Errorf("dB(FromDB(%g)) = %g, want %g", c.db, back, c.db)
		}
	}
}

func TestDBmWattsRoundTrip(t *testing.T) {
	watts := DBmToWatts(0)
	if math.Abs(watts-1e-3) > 1e-15 {
		t.Errorf("DBmToWatts(0) = %g W, want 0.001 W", watts)
	}
	back, err := WattsToDBm(watts)
	if err != nil {
		t.Fatalf("WattsToDBm(%g) returned error: %v", watts, err)
	}
	if math.Abs(back) > 1e-12 {
		t.Errorf("round trip 0 dBm -> W -> dBm = %g, want 0", back)
	}
}

func TestGainLinearToDBRejectsNonPositive(t *testing.T) {
	for _, g := range []float64{0, -1, -0.5} {
		if _, err := GainLinearToDB(g); err == nil {
			t.Errorf("GainLinearToDB(%g) = nil error, want rejection", g)
		}
	}
	got, err := GainLinearToDB(100)
	if err != nil {
		t.Fatalf("GainLinearToDB(100) returned error: %v", err)
	}
	if math.Abs(got-20) > 1e-12 {
		t.Errorf("GainLinearToDB(100) = %g dB, want 20 dB", got)
	}
}

func TestSpeedOfLightPinned(t *testing.T) {
	if SpeedOfLight != 2.99792458e8 {
		t.Errorf("SpeedOfLight = %g, want 2.99792458e8", SpeedOfLight)
	}
}

func TestBoltzmannPinned(t *testing.T) {
	if Boltzmann != 1.380649e-23 {
		t.Errorf("Boltzmann = %g, want 1.380649e-23", Boltzmann)
	}
}
