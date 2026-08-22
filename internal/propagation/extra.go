package propagation

import "math"

// ExtraLoss models additional attenuation that is not part of the
// free-space path itself: polarisation mismatch, antenna pointing
// error, feeder or connector loss. It is expressed in decibels and is
// subtracted exactly once from the received power.
type ExtraLoss struct {
	DB float64
}

// NewExtraLoss constructs an extra loss from a decibel value. The value
// must be non-negative.
func NewExtraLoss(db float64) (ExtraLoss, error) {
	if db < 0 {
		return ExtraLoss{}, errNegativeExtraLoss(db)
	}
	return ExtraLoss{DB: db}, nil
}

// Validate returns an error when the extra loss is negative.
func (l ExtraLoss) Validate() error {
	if l.DB < 0 {
		return errNegativeExtraLoss(l.DB)
	}
	return nil
}

// Apply subtracts the extra loss from a received power level in dBm.
func (l ExtraLoss) Apply(powerDBm float64) float64 {
	return powerDBm - l.DB
}

// Inverse returns the extra loss as a linear power ratio. Provided for
// callers that stay entirely in the linear domain; the decibel path is
// the primary one.
func (l ExtraLoss) Inverse() float64 {
	return math.Pow(10, -l.DB/10)
}
