package propagation

import "math"

type ExtraLoss struct {
	DB float64
}

func NewExtraLoss(db float64) (ExtraLoss, error) {
	if db < 0 {
		return ExtraLoss{}, errNegativeExtraLoss(db)
	}
	return ExtraLoss{DB: db}, nil
}

func (l ExtraLoss) Validate() error {
	if l.DB < 0 {
		return errNegativeExtraLoss(l.DB)
	}
	return nil
}

func (l ExtraLoss) Apply(powerDBm float64) float64 {
	return powerDBm - l.DB
}

func (l ExtraLoss) Inverse() float64 {
	return math.Pow(10, -l.DB/10)
}
