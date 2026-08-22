package propagation

import (
	"fmt"

	"friis-link/internal/model"
)

// EffectiveAreaM returns the effective aperture of a receive antenna:
//
//	Ae = lambda^2 * G / (4*pi)
//
// with G the linear gain and lambda the wavelength. The Friis equation
// can be read as Pr = S * Ae, where S is the incident power density.
func EffectiveAreaM(gainLinear, frequencyHz float64) (float64, error) {
	if err := model.RequireGainLinear("rx_gain", gainLinear); err != nil {
		return 0, err
	}
	lambda, err := Wavelength(frequencyHz)
	if err != nil {
		return 0, err
	}
	return lambda * lambda * gainLinear / (4 * model.Pi), nil
}

// PowerDensityWPerM2 returns the power flux density at distance d from
// an isotropic radiator of equivalent radiated power eirpW:
//
//	S = EIRP / (4 * pi * d^2)
func PowerDensityWPerM2(eirpW, distanceM float64) (float64, error) {
	if err := model.RequireNonNegative("eirp", eirpW); err != nil {
		return 0, err
	}
	if err := model.RequirePositive("distance", distanceM); err != nil {
		return 0, err
	}
	return eirpW / (4 * model.Pi * distanceM * distanceM), nil
}

// ReceivedFromPowerDensity combines an incident power density with an
// effective aperture into the captured power:
//
//	Pr = S * Ae
func ReceivedFromPowerDensity(sWPerM2, aeM2 float64) (float64, error) {
	if sWPerM2 < 0 {
		return 0, fmt.Errorf("power density %g is negative", sWPerM2)
	}
	if aeM2 < 0 {
		return 0, fmt.Errorf("effective area %g is negative", aeM2)
	}
	return sWPerM2 * aeM2, nil
}

// DecomposedLink is the three-step view of the Friis equation.
type DecomposedLink struct {
	// PowerDensity is EIRP/(4*pi*d^2) at the receiver.
	PowerDensity float64
	// EffectiveArea is lambda^2*Gr/(4*pi).
	EffectiveArea float64
	// ReceivedPower is the product of the two.
	ReceivedPower float64
}

// DecomposeFriis evaluates the link as S * Ae. It must agree with
// ReceivedPowerLinear; the test pins the agreement so the decomposed
// view cannot drift from the direct formula.
func DecomposeFriis(ptW, gtLinear, grLinear, distanceM, frequencyHz float64) (DecomposedLink, error) {
	eirpW, err := EIRPWatts(ptW, gtLinear)
	if err != nil {
		return DecomposedLink{}, err
	}
	s, err := PowerDensityWPerM2(eirpW, distanceM)
	if err != nil {
		return DecomposedLink{}, err
	}
	ae, err := EffectiveAreaM(grLinear, frequencyHz)
	if err != nil {
		return DecomposedLink{}, err
	}
	pr, err := ReceivedFromPowerDensity(s, ae)
	if err != nil {
		return DecomposedLink{}, err
	}
	return DecomposedLink{
		PowerDensity:  s,
		EffectiveArea: ae,
		ReceivedPower: pr,
	}, nil
}
