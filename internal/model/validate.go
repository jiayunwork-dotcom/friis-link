package model

import (
	"fmt"
	"math"
)

func RequirePositive(name string, value float64) error {
	if value <= 0 {
		return fmt.Errorf("%s: %g must be strictly positive", name, value)
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%s: %g is not a finite number", name, value)
	}
	return nil
}

func RequireNonNegative(name string, value float64) error {
	if value < 0 {
		return fmt.Errorf("%s: %g must be non-negative", name, value)
	}
	if math.IsNaN(value) {
		return fmt.Errorf("%s: %g is not a number", name, value)
	}
	return nil
}

func RequireGainLinear(name string, gainLinear float64) error {
	if gainLinear < 0 {
		return fmt.Errorf("%s: linear gain %g is negative", name, gainLinear)
	}
	if gainLinear == 0 {
		return fmt.Errorf("%s: linear gain %g must be positive", name, gainLinear)
	}
	if math.IsNaN(gainLinear) || math.IsInf(gainLinear, 0) {
		return fmt.Errorf("%s: linear gain %g is not a finite number", name, gainLinear)
	}
	return nil
}
