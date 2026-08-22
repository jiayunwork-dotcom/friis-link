package budget

import "errors"

// errNilConfig guards against a nil pointer reaching the kernel.
var errNilConfig = errors.New("budget: nil config")

// errCrossCheckFailed is returned when the linear and the decibel
// formulation of the received power disagree by more than the tolerance.
// It protects the invariant that both code paths share one conversion.
var errCrossCheckFailed = errors.New("budget: linear and decibel received power disagree")
