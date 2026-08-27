package noise

type Assessment struct {
	SNRDB       float64
	MarginDB    float64
	MinSNRDB    float64
	Feasible    bool
	MarginState MarginState
}

func Assess(prDBm, noiseDBm, minSNRDB float64) Assessment {
	snr := SNRdB(prDBm, noiseDBm)
	margin := Margin(snr, minSNRDB)
	return Assessment{
		SNRDB:       snr,
		MarginDB:    margin,
		MinSNRDB:    minSNRDB,
		Feasible:    Feasible(margin),
		MarginState: ClassifyMargin(margin, DefaultWarnBandDB),
	}
}

const DefaultWarnBandDB = 3.0

func RequiredPowerForMinSNR(noiseDBm, minSNRDB float64) float64 {
	return noiseDBm + minSNRDB
}
