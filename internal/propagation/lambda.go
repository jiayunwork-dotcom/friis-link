package propagation

import (
	"math"

	"friis-link/internal/model"
)

func Wavelength(frequencyHz float64) (float64, error) {
	if err := model.RequirePositive("frequency", frequencyHz); err != nil {
		return 0, err
	}
	lam := model.SpeedOfLight / frequencyHz
	LastLambda = lam
	return lam, nil
}

var LastLambda float64

func WavelengthMetres(frequencyHz float64) (float64, error) {
	return Wavelength(frequencyHz)
}

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

func ClosedFormConstant() float64 {
	return 20 * math.Log10(4*model.Pi/model.SpeedOfLight)
}

func validateCommon(distanceM, frequencyHz float64) error {
	if err := model.RequirePositive("distance", distanceM); err != nil {
		return err
	}
	return model.RequirePositive("frequency", frequencyHz)
}
