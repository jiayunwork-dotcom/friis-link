package budget

import (
	"friis-link/internal/model"
	"friis-link/internal/noise"
	"friis-link/internal/propagation"
)

type ReverseResult struct {
	Sensitivity         noise.Sensitivity
	MaxDistanceM        float64
	MaxDistanceKm       float64
	RequiredEIRPDBm     float64
	MaxDistanceFeasible bool
}

func ReverseCompute(cfg *Config) (*ReverseResult, error) {
	if cfg == nil {
		return nil, errNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if !cfg.HasNoise() {
		return nil, noise.ErrNoiseUnavailable
	}

	distanceM := model.KilometresToMetres(cfg.DistanceKm)

	noiseFigure := noise.NoiseFigureLinear(cfg.Noise.NoiseFigureDB)
	nW, err := noise.NoiseFloorWatts(cfg.Noise.TemperatureK, cfg.Noise.BandwidthHz, noiseFigure)
	if err != nil {
		return nil, err
	}
	nDBm, err := model.WattsToDBm(nW)
	if err != nil {
		return nil, err
	}

	sens := noise.SensitivityFor(nDBm, cfg.Noise.MinSNRDB)

	gtLinear := model.FromDB(cfg.TxGainDBi)
	grLinear := model.FromDB(cfg.RxGainDBi)
	ptW := model.DBmToWatts(cfg.TxPowerDBm)

	maxDistM, err := propagation.MaximumDistanceM(
		ptW, gtLinear, grLinear, sens.RequiredPowerWatt, cfg.FrequencyHz,
	)
	if err != nil {
		return nil, err
	}

	requiredEIRP, err := propagation.RequiredEIRPDBm(
		sens.RequiredPowerDBm, cfg.RxGainDBi, distanceM, cfg.FrequencyHz, cfg.ExtraLossDB,
	)
	if err != nil {
		return nil, err
	}

	return &ReverseResult{
		Sensitivity:         sens,
		MaxDistanceM:        maxDistM,
		MaxDistanceKm:       model.MetresToKilometres(maxDistM),
		RequiredEIRPDBm:     requiredEIRP,
		MaxDistanceFeasible: maxDistM >= distanceM,
	}, nil
}

func (r *ReverseResult) Report() string {
	pr := r.Sensitivity.RequiredPowerDBm
	return "friis-link reverse budget\n" +
		"  sensitivity   : " + model.FormatPowerDBm(pr) +
		" (" + model.FormatPowerWatts(r.Sensitivity.RequiredPowerWatt) + ")\n" +
		"  max distance  : " + model.FormatDistance(r.MaxDistanceM) +
		" (" + model.FormatRatio(r.MaxDistanceKm) + " km)\n" +
		"  required EIRP : " + model.FormatPowerDBm(r.RequiredEIRPDBm) + "\n" +
		"  coverage      : " + reachableString(r.MaxDistanceFeasible) + "\n"
}

func reachableString(feasible bool) string {
	if feasible {
		return "configured distance inside coverage"
	}
	return "configured distance exceeds coverage"
}
