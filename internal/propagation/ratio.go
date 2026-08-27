package propagation

import (
	"fmt"
	"math"
)

func ratioToDB(ratio float64) (float64, error) {
	if ratio < 0 {
		return 0, fmt.Errorf("path loss ratio %g is negative", ratio)
	}
	if ratio == 0 {
		return math.Inf(-1), nil
	}
	return 10 * math.Log10(ratio), nil
}
