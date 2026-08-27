package model

import (
	"fmt"
	"math"
)

type ValidationError struct {
	Field string
	Value float64
}

func (e *ValidationError) Error() string {
	return "invalid input field " + e.Field + " (got " + formatFloat(e.Value) + ")"
}

func formatFloat(v float64) string {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return fmt.Sprintf("%g", v)
	}
	return fmt.Sprintf("%g", v)
}
