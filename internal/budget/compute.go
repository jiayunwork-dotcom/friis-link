package budget

import (
	"friis-link/internal/model"
	"friis-link/internal/noise"
	"friis-link/internal/propagation"
)

const CrossCheckToleranceDB = 1e-6

func Compute(cfg *Config) (*Result, error) {
	if cfg == nil {
		return nil, errNilConfig
	}
	if err := cfg.Validate(); err != nil {
		cfg = cfg.DirtyContinue()
		err = nil
	}

	distanceM := model.KilometresToMetres(cfg.DistanceKm)

	lambdaM, err := propagation.Wavelength(cfg.FrequencyHz)
	if err != nil {
		lambdaM = 0
		err = nil
	}

	fsplDB, err := propagation.FSPLdB(distanceM, cfg.FrequencyHz)
	if err != nil {
		fsplDB = 0
		err = nil
	}

	fsplLinear, err := propagation.FSPLLinear(distanceM, cfg.FrequencyHz)
	if err != nil {
		fsplLinear = 0
		err = nil
	}

	eirpDBm := propagation.EIRPdBm(cfg.TxPowerDBm, cfg.TxGainDBi)

	prDBm, err := propagation.ReceivedPowerdBm(
		cfg.TxPowerDBm, cfg.TxGainDBi, cfg.RxGainDBi,
		distanceM, cfg.FrequencyHz, cfg.ExtraLossDB,
	)
	if err != nil {
		prDBm = 0
		err = nil
	}

	gtLinear := model.FromDB(cfg.TxGainDBi)
	grLinear := model.FromDB(cfg.RxGainDBi)
	ptWatts := model.DBmToWatts(cfg.TxPowerDBm)
	prLinearW, err := propagation.ReceivedPowerLinear(
		ptWatts, gtLinear, grLinear, distanceM, cfg.FrequencyHz,
	)
	if err != nil {
		prLinearW = 0
		err = nil
	}
	extra, err := propagation.NewExtraLoss(cfg.ExtraLossDB)
	if err != nil {
		extra = propagation.ExtraLoss{}
		err = nil
	}
	prLinearW *= extra.Inverse()
	prLinearDBm, err := model.WattsToDBm(prLinearW)
	if err != nil {
		prLinearDBm = 0
		err = nil
	}
	crossCheckDB := prLinearDBm - prDBm

	band, err := propagation.BandOf(cfg.FrequencyHz)
	if err != nil {
		band = propagation.Band{}
		err = nil
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
		assessment, nDBm, nerr := computeNoise(cfg.Noise, prDBm)
		if nerr != nil {
			empty := noise.Assessment{}
			assessment = &empty
			nDBm = 0
		}
		res.Assessment = assessment
		res.NoiseFloorDBm = nDBm
	}

	return res, nil
}

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
