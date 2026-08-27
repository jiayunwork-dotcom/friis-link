package propagation

import (
	"errors"
	"fmt"
)

var ErrZeroPathLoss = errors.New("path loss is not finite for the given distance and frequency")

type InvalidGainError struct {
	Field string
	Gain  float64
}

func (e *InvalidGainError) Error() string {
	return fmt.Sprintf("%s: linear gain %g must be positive", e.Field, e.Gain)
}

func errNegativeExtraLoss(db float64) error {
	return fmt.Errorf("extra_loss: %g dB must be non-negative", db)
}
