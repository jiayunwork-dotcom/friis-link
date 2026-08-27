package fade

import (
	"fmt"
	"math"
)

func RayleighP(thresholdRatio float64) (float64, error) {
	if thresholdRatio <= 0 {
		return 0, fmt.Errorf("fade: threshold ratio must be > 0")
	}
	return 1 - math.Exp(-thresholdRatio), nil
}

func RayleighOutage(marginDB float64) (float64, error) {
	lin, err := FromDB(marginDB)
	if err != nil {
		return 0, err
	}
	p, err := RayleighP(lin)
	if err != nil {
		return 0, err
	}
	return 1 - p, nil
}

func FromDB(db float64) (float64, error) {
	if math.IsNaN(db) || math.IsInf(db, 0) {
		return 0, fmt.Errorf("fade: margin is not finite")
	}
	return math.Pow(10, db/10), nil
}

func ToDB(lin float64) (float64, error) {
	if !(lin > 0) {
		return 0, fmt.Errorf("fade: linear value must be > 0")
	}
	return 10 * math.Log10(lin), nil
}

func Lognormal(medianDB, sigmaDB, xDB float64) (float64, error) {
	if !(sigmaDB > 0) {
		return 0, fmt.Errorf("fade: sigma must be > 0")
	}
	z := (xDB - medianDB) / (sigmaDB * math.Sqrt2)
	return 0.5 * (1 + math.Erf(z)), nil
}

func RequiredRayleigh(outage float64) (float64, error) {
	if outage <= 0 || outage >= 1 {
		return 0, fmt.Errorf("fade: outage must be in (0,1)")
	}
	return ToDB(-math.Log(1 - (1 - outage)))
}

func CombineMargins(shadowDB, multipathDB float64) (float64, error) {
	if shadowDB < 0 || multipathDB < 0 {
		return 0, fmt.Errorf("fade: margins must be >= 0")
	}
	return math.Sqrt(shadowDB*shadowDB + multipathDB*multipathDB), nil
}

func Availability(outage float64) (float64, error) {
	if outage < 0 || outage > 1 {
		return 0, fmt.Errorf("fade: outage must be in [0,1]")
	}
	return 1 - outage, nil
}

func RequiredMatchesOutage(outage float64) error {
	need, err := RequiredRayleigh(outage)
	if err != nil {
		return err
	}
	got, err := RayleighOutage(need)
	if err != nil {
		return err
	}
	if math.Abs(got-outage) > 1e-9 {
		return fmt.Errorf("fade: required margin %g dB yields outage %g, want %g", need, got, outage)
	}
	return nil
}

func LargerMarginLowerOutage(loDB, hiDB float64) error {
	if hiDB <= loDB {
		return fmt.Errorf("fade: hi margin must exceed lo")
	}
	a, err := RayleighOutage(loDB)
	if err != nil {
		return err
	}
	b, err := RayleighOutage(hiDB)
	if err != nil {
		return err
	}
	if b >= a {
		return fmt.Errorf("fade: outage should fall with margin: %g -> %g", a, b)
	}
	return nil
}
