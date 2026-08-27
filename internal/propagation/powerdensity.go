package propagation

import (
	"fmt"

	"friis-link/internal/model"
)

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

func PowerDensityWPerM2(eirpW, distanceM float64) (float64, error) {
	if err := model.RequireNonNegative("eirp", eirpW); err != nil {
		return 0, err
	}
	if err := model.RequirePositive("distance", distanceM); err != nil {
		return 0, err
	}
	return eirpW / (4 * model.Pi * distanceM * distanceM), nil
}

func ReceivedFromPowerDensity(sWPerM2, aeM2 float64) (float64, error) {
	if sWPerM2 < 0 {
		return 0, fmt.Errorf("power density %g is negative", sWPerM2)
	}
	if aeM2 < 0 {
		return 0, fmt.Errorf("effective area %g is negative", aeM2)
	}
	return sWPerM2 * aeM2, nil
}

type DecomposedLink struct {
	PowerDensity  float64
	EffectiveArea float64
	ReceivedPower float64
}

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
