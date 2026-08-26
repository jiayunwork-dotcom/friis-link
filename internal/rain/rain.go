package rain

import (
	"fmt"
	"math"
)

func Specific(k, alpha, rateMMH float64) (float64, error) {
	if !(k > 0) {
		return 0, fmt.Errorf("rain: k must be > 0")
	}
	if alpha <= 0 {
		return 0, fmt.Errorf("rain: alpha must be > 0")
	}
	if rateMMH < 0 {
		return 0, fmt.Errorf("rain: rain rate must be >= 0")
	}
	if rateMMH == 0 {
		return 0, nil
	}
	return k * math.Pow(rateMMH, alpha), nil
}

func Path(k, alpha, rateMMH, lengthKm float64) (float64, error) {
	if !(lengthKm > 0) {
		return 0, fmt.Errorf("rain: path length must be > 0")
	}
	g, err := Specific(k, alpha, rateMMH)
	if err != nil {
		return 0, err
	}
	return g * lengthKm, nil
}

func EffectiveLength(lengthKm, reduction float64) (float64, error) {
	if !(lengthKm > 0) {
		return 0, fmt.Errorf("rain: path length must be > 0")
	}
	if reduction <= 0 || reduction > 1 {
		return 0, fmt.Errorf("rain: reduction factor must be in (0, 1]")
	}
	return lengthKm * reduction, nil
}

func PathEffective(k, alpha, rateMMH, lengthKm, reduction float64) (float64, error) {
	le, err := EffectiveLength(lengthKm, reduction)
	if err != nil {
		return 0, err
	}
	return Path(k, alpha, rateMMH, le)
}

func DoubleRate(k, alpha, rateMMH float64) (float64, float64, error) {
	a, err := Specific(k, alpha, rateMMH)
	if err != nil {
		return 0, 0, err
	}
	b, err := Specific(k, alpha, 2*rateMMH)
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

func LinearToDB(linear float64) (float64, error) {
	if !(linear > 0) {
		return 0, fmt.Errorf("rain: linear factor must be > 0")
	}
	return 10 * math.Log10(linear), nil
}

type Coeff struct {
	K     float64
	Alpha float64
}

func ITUCoeff(freqGHz float64, horizontal bool) (Coeff, error) {
	if !(freqGHz > 0) {
		return Coeff{}, fmt.Errorf("rain: frequency must be > 0")
	}
	k := 0.0001 * math.Pow(freqGHz, 1.6)
	a := 1.0
	if freqGHz >= 10 {
		k = 0.01 * math.Pow(freqGHz/10, 1.2)
		a = 1.2
	}
	if freqGHz >= 20 {
		k = 0.075 * math.Pow(freqGHz/20, 0.8)
		a = 1.05
	}
	if !horizontal {
		k *= 0.85
		a *= 0.95
	}
	if k <= 0 || a <= 0 {
		return Coeff{}, fmt.Errorf("rain: degenerate ITU coefficients")
	}
	return Coeff{K: k, Alpha: a}, nil
}

func HigherBandMoreRain(rateMMH, lengthKm float64) error {
	lo, err := ITUCoeff(8, true)
	if err != nil {
		return err
	}
	hi, err := ITUCoeff(20, true)
	if err != nil {
		return err
	}
	a, err := Path(lo.K, lo.Alpha, rateMMH, lengthKm)
	if err != nil {
		return err
	}
	b, err := Path(hi.K, hi.Alpha, rateMMH, lengthKm)
	if err != nil {
		return err
	}
	if b <= a {
		return fmt.Errorf("rain: 20 GHz should exceed 8 GHz path rain: %g vs %g", b, a)
	}
	return nil
}

func ReductionShrinksPath(k, alpha, rateMMH, lengthKm float64) error {
	full, err := Path(k, alpha, rateMMH, lengthKm)
	if err != nil {
		return err
	}
	red, err := PathEffective(k, alpha, rateMMH, lengthKm, 0.5)
	if err != nil {
		return err
	}
	if math.Abs(2*red-full) > 1e-9 {
		return fmt.Errorf("rain: half reduction should halve path: %g vs %g", red, full)
	}
	return nil
}
