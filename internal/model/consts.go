// Package model holds the physical constants, unit helpers and dB
// conversion routines shared by every part of the link budget.
//
// All propagation math in this repository goes through the two pinned
// constants declared here. The speed of light and the Boltzmann constant
// are fixed once and must never drift between packages.
package model

import "math"

// Physical constants, pinned for the whole repository.
const (
	// SpeedOfLight is the vacuum speed of light in metres per second.
	SpeedOfLight = 2.99792458e8

	// Boltzmann is the Boltzmann constant in joules per kelvin.
	Boltzmann = 1.380649e-23
)

// Pi is exported so callers do not mix half-written literals of pi
// across packages. Everything trigonometric uses the same value.
const Pi = math.Pi
