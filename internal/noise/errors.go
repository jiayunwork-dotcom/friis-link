package noise

import (
	"errors"
	"fmt"
)

var ErrNoiseUnavailable = errors.New("noise parameters were not provided")

var ErrUndefinedSNR = errors.New("SNR is undefined for the given noise input")

func WrapNoiseError(stage string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("noise[%s]: %w", stage, err)
}
