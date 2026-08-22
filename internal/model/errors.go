package model

import (
	"fmt"
	"math"
)

// ValidationError wraps a failed input check together with the name of
// the field that failed. The name is the stable key: tests assert on
// the field name, never on the full message wording.
type ValidationError struct {
	Field string
	Value float64
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return "invalid input field " + e.Field + " (got " + formatFloat(e.Value) + ")"
}

func formatFloat(v float64) string {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return fmt.Sprintf("%g", v)
	}
	return fmt.Sprintf("%g", v)
}
