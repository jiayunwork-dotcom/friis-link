package budget

import "errors"

var errNilConfig = errors.New("budget: nil config")

var errCrossCheckFailed = errors.New("budget: linear and decibel received power disagree")
