package slant

import (
	"fmt"
	"math"

	"friis-link/internal/rain"
)

func ElevationRad(elevDeg float64) (float64, error) {
	if elevDeg <= 0 || elevDeg > 90 {
		return 0, fmt.Errorf("slant: elevation must be in (0, 90]")
	}
	return elevDeg * math.Pi / 180, nil
}

func SlantRange(groundKm, elevDeg float64) (float64, error) {
	if !(groundKm > 0) {
		return 0, fmt.Errorf("slant: ground distance must be > 0")
	}
	el, err := ElevationRad(elevDeg)
	if err != nil {
		return 0, err
	}
	s := math.Sin(el)
	if s <= 0 {
		return 0, fmt.Errorf("slant: sine of elevation is not positive")
	}
	return groundKm / s, nil
}

func GroundFromSlant(slantKm, elevDeg float64) (float64, error) {
	if !(slantKm > 0) {
		return 0, fmt.Errorf("slant: slant range must be > 0")
	}
	el, err := ElevationRad(elevDeg)
	if err != nil {
		return 0, err
	}
	return slantKm * math.Sin(el), nil
}

func RainEffective(groundKm, elevDeg, reduction float64) (float64, error) {
	sl, err := SlantRange(groundKm, elevDeg)
	if err != nil {
		return 0, err
	}
	return rain.EffectiveLength(sl, reduction)
}

func RainPath(k, alpha, rateMMH, groundKm, elevDeg, reduction float64) (float64, error) {
	le, err := RainEffective(groundKm, elevDeg, reduction)
	if err != nil {
		return 0, err
	}
	return rain.Path(k, alpha, rateMMH, le)
}

func HigherElevationShortensRain(k, alpha, rateMMH, groundKm, lowDeg, highDeg, reduction float64) error {
	if highDeg <= lowDeg {
		return fmt.Errorf("slant: high elevation must exceed low")
	}
	a, err := RainPath(k, alpha, rateMMH, groundKm, lowDeg, reduction)
	if err != nil {
		return err
	}
	b, err := RainPath(k, alpha, rateMMH, groundKm, highDeg, reduction)
	if err != nil {
		return err
	}
	if b >= a {
		return fmt.Errorf("slant: higher elevation should cut rain path %g -> %g", a, b)
	}
	return nil
}

func SlantExceedsGround(groundKm, elevDeg float64) error {
	sl, err := SlantRange(groundKm, elevDeg)
	if err != nil {
		return err
	}
	if elevDeg < 90 && sl <= groundKm {
		return fmt.Errorf("slant: range %g should exceed ground %g", sl, groundKm)
	}
	if elevDeg == 90 && math.Abs(sl-groundKm) > 1e-12 {
		return fmt.Errorf("slant: zenith range should equal height")
	}
	return nil
}

func RoundTripGround(groundKm, elevDeg float64) error {
	sl, err := SlantRange(groundKm, elevDeg)
	if err != nil {
		return err
	}
	back, err := GroundFromSlant(sl, elevDeg)
	if err != nil {
		return err
	}
	if math.Abs(back-groundKm) > 1e-9 {
		return fmt.Errorf("slant: ground round-trip %g != %g", back, groundKm)
	}
	return nil
}

func ZenithIsIdentity(heightKm float64) error {
	sl, err := SlantRange(heightKm, 90)
	if err != nil {
		return err
	}
	if math.Abs(sl-heightKm) > 1e-12 {
		return fmt.Errorf("slant: zenith %g != %g", sl, heightKm)
	}
	return nil
}
