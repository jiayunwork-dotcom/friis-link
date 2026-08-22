// Package noise implements the thermal noise floor and the SNR/margin
// assessment of a received signal.
//
// The noise floor is N = k*T*B*F with the Boltzmann constant k pinned
// in the model package, T the system temperature in kelvin, B the
// receiver bandwidth in hertz and F the noise figure as a linear ratio.
// Doubling the bandwidth doubles the noise floor, which lowers the SNR
// by exactly 3 dB; the tests pin that scaling.
package noise

import (
	"fmt"

	"friis-link/internal/model"
)

// NoiseFloorWatts returns N = k*T*B*F in watts. All three quantities
// must be strictly positive; a zero bandwidth is rejected so the SNR
// never becomes infinite.
func NoiseFloorWatts(temperatureK, bandwidthHz, noiseFigureLinear float64) (float64, error) {
	if err := model.RequirePositive("noise_temperature", temperatureK); err != nil {
		return 0, err
	}
	if err := model.RequirePositive("bandwidth", bandwidthHz); err != nil {
		return 0, err
	}
	if err := model.RequirePositive("noise_figure", noiseFigureLinear); err != nil {
		return 0, err
	}
	return model.Boltzmann * temperatureK * bandwidthHz * noiseFigureLinear, nil
}

// NoiseFloordBm returns the same noise floor expressed in dBm.
func NoiseFloordBm(temperatureK, bandwidthHz, noiseFigureLinear float64) (float64, error) {
	nW, err := NoiseFloorWatts(temperatureK, bandwidthHz, noiseFigureLinear)
	if err != nil {
		return 0, err
	}
	return model.WattsToDBm(nW)
}

// NoiseFigureLinear converts a noise figure from dB to a linear ratio.
// A 0 dB noise figure corresponds to a linear factor of one.
func NoiseFigureLinear(noiseFigureDB float64) float64 {
	return model.FromDB(noiseFigureDB)
}

// ReferenceNoiseFloorAt computes the noise floor for the textbook 290 K
// reference temperature; kept separate from NoiseFloorWatts so reports
// can show both the actual and the reference floor.
func ReferenceNoiseFloorAt(bandwidthHz float64) (float64, error) {
	return NoiseFloorWatts(290, bandwidthHz, 1)
}

// ValidateNoiseInputs performs the same checks as NoiseFloorWatts
// without computing the product. It is used by the budget kernel to
// fail fast before any log-domain arithmetic runs.
func ValidateNoiseInputs(temperatureK, bandwidthHz, noiseFigureLinear float64) error {
	if err := model.RequirePositive("noise_temperature", temperatureK); err != nil {
		return err
	}
	if err := model.RequirePositive("bandwidth", bandwidthHz); err != nil {
		return err
	}
	return model.RequirePositive("noise_figure", noiseFigureLinear)
}

// formatNoise is a tiny helper that keeps the noise floor report
// consistent between the linear and decibel views.
func formatNoise(watts float64) (string, error) {
	dbm, err := model.WattsToDBm(watts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.6g W (%.3f dBm)", watts, dbm), nil
}
