package model

import (
	"fmt"
	"math"
)

// Common validation rules shared by propagation and noise. Every error
// message names the offending field so the CLI can print something an
// operator can act on without digging into the source.

// RequirePositive rejects a quantity that must be strictly greater than
// zero (distance, frequency, bandwidth, noise temperature, noise figure).
func RequirePositive(name string, value float64) error {
	if value <= 0 {
		return fmt.Errorf("%s: %g must be strictly positive", name, value)
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%s: %g is not a finite number", name, value)
	}
	return nil
}

// RequireNonNegative rejects a quantity that cannot go below zero but
// may legitimately be zero (extra loss in dB defaults to zero).
func RequireNonNegative(name string, value float64) error {
	if value < 0 {
		return fmt.Errorf("%s: %g must be non-negative", name, value)
	}
	if math.IsNaN(value) {
		return fmt.Errorf("%s: %g is not a number", name, value)
	}
	return nil
}

// RequireGainLinear rejects a linear gain that is negative. A linear
// gain of zero is also rejected because its decibel equivalent is not
// finite and would poison every downstream log-domain term.
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
