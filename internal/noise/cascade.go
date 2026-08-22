package noise

import (
	"fmt"

	"friis-link/internal/model"
)

// Stage is one block of a receiver chain: a linear gain and a linear
// noise figure. The order of stages matters for the cascade formula.
type Stage struct {
	GainLinear    float64
	NoiseFigure   float64
	Name          string
}

// CascadeNoiseFigure computes the total noise figure of a receiver
// chain with the Friis cascade formula:
//
//	F_total = F1 + (F2-1)/G1 + (F3-1)/(G1*G2) + ...
//
// Each noise figure and gain is linear. The result is the linear noise
// figure of the whole chain; convert with model.GainLinearToDB when a
// decibel figure is wanted.
func CascadeNoiseFigure(stages []Stage) (float64, error) {
	if len(stages) == 0 {
		return 0, fmt.Errorf("cascade: empty stage list")
	}
	total := 0.0
	productGain := 1.0
	for i, s := range stages {
		if err := model.RequirePositive("gain", s.GainLinear); err != nil {
			return 0, fmt.Errorf("cascade stage %d (%s): %w", i, s.Name, err)
		}
		if err := model.RequirePositive("noise_figure", s.NoiseFigure); err != nil {
			return 0, fmt.Errorf("cascade stage %d (%s): %w", i, s.Name, err)
		}
		if i == 0 {
			total = s.NoiseFigure
		} else {
			total += (s.NoiseFigure - 1) / productGain
		}
		productGain *= s.GainLinear
	}
	return total, nil
}

// CascadedNoiseFloorWatts computes N = k*T*B*F_cascade for a receiver
// chain whose total noise figure is F_cascade.
func CascadedNoiseFloorWatts(temperatureK, bandwidthHz float64, stages []Stage) (float64, error) {
	fc, err := CascadeNoiseFigure(stages)
	if err != nil {
		return 0, err
	}
	return NoiseFloorWatts(temperatureK, bandwidthHz, fc)
}

// LNAStages is a helper that builds a canonical three-stage chain
// (LNA, mixer, IF amplifier) with the given linear gains and noise
// figures. It exists so the report can show a named receiver chain
// instead of anonymous stages.
func LNAStages(lnaGain, lnaNF, mixerGain, mixerNF, ifGain, ifNF float64) []Stage {
	return []Stage{
		{Name: "LNA", GainLinear: lnaGain, NoiseFigure: lnaNF},
		{Name: "mixer", GainLinear: mixerGain, NoiseFigure: mixerNF},
		{Name: "IF amp", GainLinear: ifGain, NoiseFigure: ifNF},
	}
}
