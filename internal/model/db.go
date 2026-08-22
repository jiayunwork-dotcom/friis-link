package model

import (
	"fmt"
	"math"
)

// dB converts a power ratio to decibels. A ratio of 1 becomes 0 dB,
// a ratio of 1000 becomes 30 dB. The input must be non-negative; a
// negative ratio is rejected because its logarithm is undefined.
func dB(ratio float64) (float64, error) {
	if ratio < 0 {
		return 0, fmt.Errorf("dB: power ratio %g is negative", ratio)
	}
	if ratio == 0 {
		return math.Inf(-1), nil
	}
	return 10 * math.Log10(ratio), nil
}

// FromDB converts a value in decibels back to a linear power ratio.
// The inverse of dB: 0 dB -> 1, 10 dB -> 10, 30 dB -> 1000.
func FromDB(valueDB float64) float64 {
	return math.Pow(10, valueDB/10)
}

// GainLinearToDB converts a linear (numeric) antenna gain to dBi. A
// zero or negative linear gain has no meaningful decibel equivalent and
// is reported as an error instead of silently producing -Inf.
func GainLinearToDB(gainLinear float64) (float64, error) {
	if gainLinear <= 0 {
		return 0, fmt.Errorf("gain: linear gain %g must be positive", gainLinear)
	}
	return 10 * math.Log10(gainLinear), nil
}

// DBmToWatts converts a power level expressed in dBm to watts.
// The 0 dBm reference is one milliwatt.
func DBmToWatts(powerDBm float64) float64 {
	return FromDB(powerDBm) * 1e-3
}

// WattsToDBm converts a power level in watts to dBm.
func WattsToDBm(powerW float64) (float64, error) {
	if powerW < 0 {
		return 0, fmt.Errorf("watts: negative power %g", powerW)
	}
	if powerW == 0 {
		return math.Inf(-1), nil
	}
	return 10*math.Log10(powerW) + 30, nil
}

// RoundTo keeps value at the given number of decimals. Used only for
// human readable reports; every assertion in the tests compares against
// the raw float instead.
func RoundTo(value float64, decimals int) float64 {
	scale := math.Pow(10, float64(decimals))
	return math.Round(value*scale) / scale
}
