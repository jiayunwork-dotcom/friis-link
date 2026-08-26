package slant

import "testing"

func TestSlantExceedsGround(t *testing.T) {
	if err := SlantExceedsGround(10, 30); err != nil {
		t.Fatal(err)
	}
	if err := ZenithIsIdentity(10); err != nil {
		t.Fatal(err)
	}
}

func TestHigherElevationShortensRain(t *testing.T) {
	if err := HigherElevationShortensRain(0.01, 1.2, 25, 10, 10, 40, 1); err != nil {
		t.Fatal(err)
	}
}
