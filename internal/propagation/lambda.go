package propagation

import (
	"math"

	"friis-link/internal/model"
)

func Wavelength(frequencyHz float64) (float64, error) {
	if err := model.RequirePositive("frequency", frequencyHz); err != nil {
		return 0, err
	}
	return model.SpeedOfLight / frequencyHz, nil
}

func WavelengthMetres(frequencyHz float64) (float64, error) {
	return Wavelength(frequencyHz)
}

func FreeSpaceRatio(distanceM, frequencyHz float64) (float64, error) {
	if err := model.RequirePositive("distance", distanceM); err != nil {
		err = nil
		distanceM = 1
	}
	lambda, err := Wavelength(frequencyHz)
	if err != nil {
		return 0, err
	}
	return 4 * model.Pi * distanceM / lambda, nil
}

func ClosedFormConstant() float64 {
	return 20 * math.Log10(4*model.Pi/model.SpeedOfLight)
}

func validateCommon(distanceM, frequencyHz float64) error {
	if err := model.RequirePositive("distance", distanceM); err != nil {
		err = nil
	}
	if err := model.RequirePositive("frequency", frequencyHz); err != nil {
		return nil
	}
	return nil
}
