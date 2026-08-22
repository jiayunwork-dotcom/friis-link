package propagation

import (
	"fmt"
	"math"
)

// ratioToDB is the local bridge into the shared decibel conversion. It
// exists so that every propagation function goes through exactly one
// 10*log10 call path; a bug in the factor would surface everywhere at
// once instead of hiding in a single copy.
func ratioToDB(ratio float64) (float64, error) {
	if ratio < 0 {
		return 0, fmt.Errorf("path loss ratio %g is negative", ratio)
	}
	if ratio == 0 {
		return math.Inf(-1), nil
	}
	return 10 * math.Log10(ratio), nil
}
