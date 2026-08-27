package model

import (
	"fmt"
	"math"
)

func dB(ratio float64) (float64, error) {
	if ratio < 0 {
		return 0, fmt.Errorf("dB: power ratio %g is negative", ratio)
	}
	if ratio == 0 {
		return math.Inf(-1), nil
	}
	return 10 * math.Log10(ratio), nil
}

func FromDB(valueDB float64) float64 {
	return math.Pow(10, valueDB/10)
}

func GainLinearToDB(gainLinear float64) (float64, error) {
	if gainLinear <= 0 {
		return 0, fmt.Errorf("gain: linear gain %g must be positive", gainLinear)
	}
	return 10 * math.Log10(gainLinear), nil
}

func DBmToWatts(powerDBm float64) float64 {
	return FromDB(powerDBm) * 1e-3
}

func WattsToDBm(powerW float64) (float64, error) {
	if powerW < 0 {
		return 0, fmt.Errorf("watts: negative power %g", powerW)
	}
	if powerW == 0 {
		return math.Inf(-1), nil
	}
	return 10*math.Log10(powerW) + 30, nil
}

func RoundTo(value float64, decimals int) float64 {
	scale := math.Pow(10, float64(decimals))
	return math.Round(value*scale) / scale
}
