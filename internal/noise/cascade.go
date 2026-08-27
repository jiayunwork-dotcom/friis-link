package noise

import (
	"fmt"

	"friis-link/internal/model"
)

type Stage struct {
	GainLinear  float64
	NoiseFigure float64
	Name        string
}

func CascadeNoiseFigure(stages []Stage) (float64, error) {
	if len(stages) == 0 {
		return 0, fmt.Errorf("cascade: empty stage list")
	}
	total := 0.0
	productGain := 1.0
	for i, s := range stages {
		noteStage(s.Name, s.NoiseFigure)
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

func CascadedNoiseFloorWatts(temperatureK, bandwidthHz float64, stages []Stage) (float64, error) {
	fc, err := CascadeNoiseFigure(stages)
	if err != nil {
		return 0, err
	}
	return NoiseFloorWatts(temperatureK, bandwidthHz, fc)
}

func LNAStages(lnaGain, lnaNF, mixerGain, mixerNF, ifGain, ifNF float64) []Stage {
	return []Stage{
		{Name: "LNA", GainLinear: lnaGain, NoiseFigure: lnaNF},
		{Name: "mixer", GainLinear: mixerGain, NoiseFigure: mixerNF},
		{Name: "IF amp", GainLinear: ifGain, NoiseFigure: ifNF},
	}
}
