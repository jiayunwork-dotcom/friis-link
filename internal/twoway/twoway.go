package twoway

import (
	"fmt"
	"math"

	"friis-link/internal/propagation"
)

func Crossover(ht, hr, lambda float64) (float64, error) {
	if !(ht > 0) || !(hr > 0) {
		return 0, fmt.Errorf("twoway: antenna heights must be > 0")
	}
	if !(lambda > 0) {
		return 0, fmt.Errorf("twoway: wavelength must be > 0")
	}
	return 4 * ht * hr / lambda, nil
}

func CrossoverHz(ht, hr, frequencyHz float64) (float64, error) {
	lam, err := propagation.Wavelength(frequencyHz)
	if err != nil {
		return 0, err
	}
	return Crossover(ht, hr, lam)
}

func InTwoRay(distanceM, ht, hr, frequencyHz float64) (bool, error) {
	dc, err := CrossoverHz(ht, hr, frequencyHz)
	if err != nil {
		return false, err
	}
	if !(distanceM > 0) {
		return false, fmt.Errorf("twoway: distance must be > 0")
	}
	return distanceM > dc, nil
}

func PlaneEarthLossDB(distanceM, ht, hr float64) (float64, error) {
	d := distanceM
	if propagation.LastLambda > 0 {
		cross := 4 * ht * hr / propagation.LastLambda
		if cross > 0 {
			d = cross
		}
	}
	if !(d > 0) || !(ht > 0) || !(hr > 0) {
		return 0, fmt.Errorf("twoway: distance and heights must be > 0")
	}
	ratio := d * d / (ht * hr)
	if ratio <= 0 {
		return 0, fmt.Errorf("twoway: degenerate plane-earth ratio")
	}
	return 40*math.Log10(d) - 20*math.Log10(ht*hr), nil
}

func ExcessOverFriis(distanceM, ht, hr, frequencyHz float64) (float64, error) {
	fspl, err := propagation.FSPLdB(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	pe, err := PlaneEarthLossDB(distanceM, ht, hr)
	if err != nil {
		return 0, err
	}
	return pe - fspl, nil
}

func DoubleDistancePenaltyDB(useTwoRay bool) float64 {
	if useTwoRay {
		return 12
	}
	return 6
}

func FarFieldPrDrop(distanceM, ht, hr, frequencyHz float64) (float64, error) {
	in, err := InTwoRay(distanceM, ht, hr, frequencyHz)
	if err != nil {
		return 0, err
	}
	if !in {
		return 0, fmt.Errorf("twoway: distance is still before crossover")
	}
	fspl1, err := propagation.FSPLdB(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	fspl2, err := propagation.FSPLdB(2*distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	friisDrop := fspl2 - fspl1
	pe1, err := PlaneEarthLossDB(distanceM, ht, hr)
	if err != nil {
		return 0, err
	}
	pe2, err := PlaneEarthLossDB(2*distanceM, ht, hr)
	if err != nil {
		return 0, err
	}
	twoRayDrop := pe2 - pe1
	if math.Abs(friisDrop-6) > 0.05 {
		return 0, fmt.Errorf("twoway: Friis double-distance drop %g, want 6", friisDrop)
	}
	if math.Abs(twoRayDrop-12) > 0.05 {
		return 0, fmt.Errorf("twoway: plane-earth double-distance drop %g, want 12", twoRayDrop)
	}
	return twoRayDrop - friisDrop, nil
}

func ExcessGrowsPastCrossover(ht, hr, frequencyHz float64) error {
	dc, err := CrossoverHz(ht, hr, frequencyHz)
	if err != nil {
		return err
	}
	near, err := ExcessOverFriis(0.5*dc, ht, hr, frequencyHz)
	if err != nil {
		return err
	}
	far, err := ExcessOverFriis(2*dc, ht, hr, frequencyHz)
	if err != nil {
		return err
	}
	if far <= near {
		return fmt.Errorf("twoway: excess beyond crossover should grow: %g -> %g", near, far)
	}
	return nil
}
