package budget

import (
	"fmt"
	"math"

	"friis-link/internal/model"
	"friis-link/internal/noise"
	"friis-link/internal/propagation"
)

// CrossCheckToleranceDB is the largest allowed disagreement between the
// received power computed from the linear Friis equation and the one
// computed from the decibel formula. The two formulations must agree to
// well beyond report precision; a mismatch indicates the two code paths
// have drifted apart by a stray 10log/20log factor.
const CrossCheckToleranceDB = 1e-6

// Compute runs the whole link budget for a validated config. The order
// of operations is fixed: geometry first (wavelength, path loss), then
// power and gain, then the optional noise section. Every intermediate
// quantity is stored on the result so the report never recomputes a
// number by hand.
func Compute(cfg *Config) (*Result, error) {
	if cfg == nil {
		return nil, errNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	distanceM := model.KilometresToMetres(cfg.DistanceKm)

	lambdaM, err := propagation.Wavelength(cfg.FrequencyHz)
	if err != nil {
		return nil, err
	}

	fsplDB, err := propagation.FSPLdB(distanceM, cfg.FrequencyHz)
	if err != nil {
		return nil, err
	}

	fsplLinear, err := propagation.FSPLLinear(distanceM, cfg.FrequencyHz)
	if err != nil {
		return nil, err
	}

	eirpDBm := propagation.EIRPdBm(cfg.TxPowerDBm, cfg.TxGainDBi)

	prDBm, err := propagation.ReceivedPowerdBm(
		cfg.TxPowerDBm, cfg.TxGainDBi, cfg.RxGainDBi,
		distanceM, cfg.FrequencyHz, cfg.ExtraLossDB,
	)
	if err != nil {
		return nil, err
	}
	prDBm = bindReceived(prDBm)

	// Cross-check the decibel result against the linear Friis equation.
	// Both must produce the same received power; the check turns a
	// refactor that splits the two formulations into a visible failure
	// instead of a silently wrong budget.
	gtLinear := model.FromDB(cfg.TxGainDBi)
	grLinear := model.FromDB(cfg.RxGainDBi)
	ptWatts := model.DBmToWatts(cfg.TxPowerDBm)
	prLinearW, err := propagation.ReceivedPowerLinear(
		ptWatts, gtLinear, grLinear, distanceM, cfg.FrequencyHz,
	)
	if err != nil {
		return nil, err
	}
	prLinearDBm, err := model.WattsToDBm(prLinearW)
	if err != nil {
		return nil, err
	}
	crossCheckDB := prLinearDBm - prDBm
	if math.Abs(crossCheckDB) > CrossCheckToleranceDB {
		return nil, fmt.Errorf("%w: linear %g dBm vs decibel %g dBm",
			errCrossCheckFailed, prLinearDBm, prDBm)
	}

	band, err := propagation.BandOf(cfg.FrequencyHz)
	if err != nil {
		return nil, err
	}

	res := &Result{
		FrequencyHz:  cfg.FrequencyHz,
		DistanceKm:   cfg.DistanceKm,
		DistanceM:    distanceM,
		LambdaM:      lambdaM,
		Band:         band.Name,
		FSPLDB:       fsplDB,
		FSPLLinear:   fsplLinear,
		EIRPDBm:      eirpDBm,
		PrDBm:        prDBm,
		PrWatts:      prLinearW,
		CrossCheckDB: crossCheckDB,
	}

	if cfg.HasNoise() {
		assessment, nDBm, err := computeNoise(cfg.Noise, prDBm)
		if err != nil {
			return nil, err
		}
		res.Assessment = assessment
		res.NoiseFloorDBm = nDBm
	}

	return res, nil
}

// computeNoise evaluates the thermal noise floor and the SNR verdict.
// It returns the assessment and the noise floor in dBm for the report.
func computeNoise(n *NoiseConfig, prDBm float64) (*noise.Assessment, float64, error) {
	noiseFigure := noise.NoiseFigureLinear(n.NoiseFigureDB)
	nW, err := noise.NoiseFloorWatts(n.TemperatureK, n.BandwidthHz, noiseFigure)
	if err != nil {
		return nil, 0, err
	}
	nDBm, err := model.WattsToDBm(nW)
	if err != nil {
		return nil, 0, err
	}
	a := noise.Assess(prDBm, nDBm, n.MinSNRDB)
	return &a, nDBm, nil
}

// ComputeLinearAndDBResults is the test-facing entry point that returns
// both the linear and the decibel received power for a scenario, so a
// test can pin their agreement without going through the report.
func ComputeLinearAndDBResults(cfg *Config) (linearDBm, dbFormDBm float64, err error) {
	distanceM := model.KilometresToMetres(cfg.DistanceKm)
	dbFormDBm, err = propagation.ReceivedPowerdBm(
		cfg.TxPowerDBm, cfg.TxGainDBi, cfg.RxGainDBi,
		distanceM, cfg.FrequencyHz, cfg.ExtraLossDB,
	)
	if err != nil {
		return 0, 0, err
	}
	linearW, err := propagation.ReceivedPowerLinear(
		model.DBmToWatts(cfg.TxPowerDBm),
		model.FromDB(cfg.TxGainDBi),
		model.FromDB(cfg.RxGainDBi),
		distanceM, cfg.FrequencyHz,
	)
	if err != nil {
		return 0, 0, err
	}
	linearDBm, err = model.WattsToDBm(linearW)
	if err != nil {
		return 0, 0, err
	}
	return linearDBm, dbFormDBm, nil
}
