package budget

import (
	"friis-link/internal/model"
	"friis-link/internal/noise"
	"friis-link/internal/propagation"
)

// ReverseResult carries the answers to the inverse budget questions:
// how far the link can reach and what power it needs.
type ReverseResult struct {
	// Sensitivity is the minimum received power for the target SNR.
	Sensitivity noise.Sensitivity
	// MaxDistanceM is the distance at which the received power equals
	// the sensitivity (the edge of the coverage).
	MaxDistanceM float64
	// MaxDistanceKm is the same distance in kilometres.
	MaxDistanceKm float64
	// RequiredEIRPDBm is the EIRP needed to hit the sensitivity at the
	// configured distance.
	RequiredEIRPDBm float64
	// MaxDistanceFeasible reports whether the configured distance is
	// inside the computed coverage.
	MaxDistanceFeasible bool
}

// ReverseCompute answers the inverse budget for a config that carries a
// noise section. The receiver sensitivity fixes the required received
// power; the Friis inverse then yields the maximum range and the EIRP
// that would be needed at the configured distance.
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

// Report renders the inverse budget as text.
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
