// Package propagation implements the free-space radio path: wavelength,
// free-space path loss (FSPL) and the received power derived from the
// Friis transmission equation.
//
// Distances enter in metres and frequencies in hertz; the caller is
// responsible for unit scaling before calling these functions. Every
// function validates its inputs and returns an error instead of
// producing an infinite or silent wrong number.
package propagation

import (
	"math"

	"friis-link/internal/model"
)

// Wavelength returns the free-space wavelength lambda = c/f for a
// frequency in hertz. A non-positive frequency is rejected.
func Wavelength(frequencyHz float64) (float64, error) {
	if err := model.RequirePositive("frequency", frequencyHz); err != nil {
		return 0, err
	}
	return model.SpeedOfLight / frequencyHz, nil
}

// WavelengthMetres is a convenience wrapper that formats the result of
// Wavelength in metres for reports.
func WavelengthMetres(frequencyHz float64) (float64, error) {
	return Wavelength(frequencyHz)
}

// FreeSpaceRatio returns 4*pi*d/lambda for a distance in metres and a
// frequency in hertz. This is the argument of the squared FSPL term;
// keeping it as its own function lets reports show the intermediate
// quantity and lets tests target the trend independently of the square.
func FreeSpaceRatio(distanceM, frequencyHz float64) (float64, error) {
	if err := model.RequirePositive("distance", distanceM); err != nil {
		return 0, err
	}
	lambda, err := Wavelength(frequencyHz)
	if err != nil {
		return 0, err
	}
	return 4 * model.Pi * distanceM / lambda, nil
}

// ClosedFormConstant is 20*log10(4*pi/c) in dB with the distance in
// metres and frequency in hertz. It is the constant term of the closed
// form FSPL_dB = 20*log10(d) + 20*log10(f) + ClosedFormConstant and is
// exposed so the S-band example can be checked against it.
func ClosedFormConstant() float64 {
	return 20 * math.Log10(4*model.Pi/model.SpeedOfLight)
}

// validateCommon rejects distance and frequency that cannot produce a
// meaningful path loss.
func validateCommon(distanceM, frequencyHz float64) error {
	if err := model.RequirePositive("distance", distanceM); err != nil {
		return err
	}
	return model.RequirePositive("frequency", frequencyHz)
}
