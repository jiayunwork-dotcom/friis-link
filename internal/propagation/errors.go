package propagation

import (
	"errors"
	"fmt"
)

// ErrZeroPathLoss is returned when the input geometry cannot produce a
// finite path loss (a zero distance or zero frequency would divide by
// zero). The caller sees an error instead of an infinite decibel value.
var ErrZeroPathLoss = errors.New("path loss is not finite for the given distance and frequency")

// InvalidGainError reports a linear gain that has no decibel equivalent.
type InvalidGainError struct {
	Field string
	Gain  float64
}

// Error implements the error interface.
func (e *InvalidGainError) Error() string {
	return fmt.Sprintf("%s: linear gain %g must be positive", e.Field, e.Gain)
}

func errNegativeExtraLoss(db float64) error {
	return fmt.Errorf("extra_loss: %g dB must be non-negative", db)
}
