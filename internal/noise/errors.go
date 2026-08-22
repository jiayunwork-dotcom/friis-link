package noise

import (
	"errors"
	"fmt"
)

// ErrNoiseUnavailable marks an assessment that could not be produced
// because the noise section was omitted from the input. The kernel
// treats it as an optional feature: propagation numbers still print,
// SNR and margin print as "n/a".
var ErrNoiseUnavailable = errors.New("noise parameters were not provided")

// ErrUndefinedSNR is returned when the SNR cannot be evaluated, for
// example a zero bandwidth that would divide by zero.
var ErrUndefinedSNR = errors.New("SNR is undefined for the given noise input")

// WrapNoiseError annotates a noise failure with the stage that hit it
// so a long report pinpoints the failing component.
func WrapNoiseError(stage string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("noise[%s]: %w", stage, err)
}
