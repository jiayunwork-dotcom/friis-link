package gas

import (
	"fmt"
	"math"
)

type Air struct {
	PressureHPa  float64
	TemperatureK float64
	VaporDensity float64
}

func StandardAir() Air {
	return Air{PressureHPa: 1013.25, TemperatureK: 288.15, VaporDensity: 7.5}
}

func (a Air) Validate() error {
	if !(a.PressureHPa > 0) {
		return fmt.Errorf("gas: pressure must be > 0")
	}
	if !(a.TemperatureK > 0) {
		return fmt.Errorf("gas: temperature must be > 0")
	}
	if a.VaporDensity < 0 {
		return fmt.Errorf("gas: vapor density must be >= 0")
	}
	return nil
}

func OxygenSpecific(freqGHz float64, a Air) (float64, error) {
	if err := a.Validate(); err != nil {
		return 0, err
	}
	if !(freqGHz > 0) {
		return 0, fmt.Errorf("gas: frequency must be > 0")
	}
	theta := 300 / a.TemperatureK
	rp := a.PressureHPa / 1013.25
	gamma := 0.182 * rp * theta * theta * (freqGHz * freqGHz) * (0.0072*rp*theta/(math.Pow(freqGHz-60, 2)+0.3*rp*rp*theta*theta) +
		1.2e-6)
	if freqGHz < 50 {
		gamma = 0.0072 * rp * theta * theta * freqGHz * freqGHz / (1 + 0.0003*freqGHz*freqGHz)
	}
	if gamma < 0 || math.IsNaN(gamma) || math.IsInf(gamma, 0) {
		return 0, fmt.Errorf("gas: oxygen specific attenuation is not finite")
	}
	return gamma, nil
}

func VaporSpecific(freqGHz float64, a Air) (float64, error) {
	if err := a.Validate(); err != nil {
		return 0, err
	}
	if !(freqGHz > 0) {
		return 0, fmt.Errorf("gas: frequency must be > 0")
	}
	theta := 300 / a.TemperatureK
	rp := a.PressureHPa / 1013.25
	rho := a.VaporDensity
	den := math.Pow(freqGHz-22.235, 2) + 0.1*rp*rp
	gamma := 0.05 * rho * theta * theta * freqGHz * freqGHz * (1 / den)
	if gamma < 0 || math.IsNaN(gamma) || math.IsInf(gamma, 0) {
		return 0, fmt.Errorf("gas: vapor specific attenuation is not finite")
	}
	return gamma, nil
}

func Specific(freqGHz float64, a Air) (float64, error) {
	o, err := OxygenSpecific(freqGHz, a)
	if err != nil {
		return 0, err
	}
	v, err := VaporSpecific(freqGHz, a)
	if err != nil {
		return 0, err
	}
	return o + v, nil
}

func Path(freqGHz, lengthKm float64, a Air) (float64, error) {
	if !(lengthKm > 0) {
		return 0, fmt.Errorf("gas: path length must be > 0")
	}
	g, err := Specific(freqGHz, a)
	if err != nil {
		return 0, err
	}
	return g * lengthKm, nil
}

func OxygenRisesToward60(loGHz, hiGHz float64, a Air) error {
	if !(loGHz > 0) || !(hiGHz > loGHz) || hiGHz >= 60 {
		return fmt.Errorf("gas: need 0 < lo < hi < 60")
	}
	a1, err := OxygenSpecific(loGHz, a)
	if err != nil {
		return err
	}
	a2, err := OxygenSpecific(hiGHz, a)
	if err != nil {
		return err
	}
	if a2 <= a1 {
		return fmt.Errorf("gas: oxygen should rise toward 60 GHz: %g -> %g", a1, a2)
	}
	return nil
}

func ScaleWithLength(freqGHz, loKm, hiKm float64, a Air) error {
	if hiKm <= loKm {
		return fmt.Errorf("gas: longer path must exceed shorter")
	}
	x, err := Path(freqGHz, loKm, a)
	if err != nil {
		return err
	}
	y, err := Path(freqGHz, hiKm, a)
	if err != nil {
		return err
	}
	want := x * (hiKm / loKm)
	if math.Abs(y-want) > 1e-9 {
		return fmt.Errorf("gas: path should scale with length: %g vs %g", y, want)
	}
	return nil
}

func DryAirLessThanMoist(freqGHz float64) error {
	dry := StandardAir()
	dry.VaporDensity = 0
	moist := StandardAir()
	moist.VaporDensity = 12
	d, err := Specific(freqGHz, dry)
	if err != nil {
		return err
	}
	m, err := Specific(freqGHz, moist)
	if err != nil {
		return err
	}
	if m <= d {
		return fmt.Errorf("gas: moist air should exceed dry: %g vs %g", m, d)
	}
	return nil
}
