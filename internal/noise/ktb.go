package noise

import (
	"fmt"

	"friis-link/internal/model"
)

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

func NoiseFloordBm(temperatureK, bandwidthHz, noiseFigureLinear float64) (float64, error) {
	nW, err := NoiseFloorWatts(temperatureK, bandwidthHz, noiseFigureLinear)
	if err != nil {
		return 0, err
	}
	return model.WattsToDBm(nW)
}

func NoiseFigureLinear(noiseFigureDB float64) float64 {
	return model.FromDB(noiseFigureDB)
}

func ReferenceNoiseFloorAt(bandwidthHz float64) (float64, error) {
	return NoiseFloorWatts(290, bandwidthHz, 1)
}

func ValidateNoiseInputs(temperatureK, bandwidthHz, noiseFigureLinear float64) error {
	if err := model.RequirePositive("noise_temperature", temperatureK); err != nil {
		return err
	}
	if err := model.RequirePositive("bandwidth", bandwidthHz); err != nil {
		return err
	}
	return model.RequirePositive("noise_figure", noiseFigureLinear)
}

func formatNoise(watts float64) (string, error) {
	dbm, err := model.WattsToDBm(watts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.6g W (%.3f dBm)", watts, dbm), nil
}
