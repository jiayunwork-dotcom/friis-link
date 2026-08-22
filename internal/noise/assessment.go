package noise

// Assessment is the complete SNR verdict for one link: the SNR itself,
// the margin against the required threshold and whether the link passes.
// The budget kernel assembles one of these and forwards it to the
// report; tests assert on the public fields.
type Assessment struct {
	SNRDB       float64
	MarginDB    float64
	MinSNRDB    float64
	Feasible    bool
	MarginState MarginState
}

// Assess computes the SNR verdict for a received power level and a noise
// floor, both in dBm, against a required minimum SNR. The required SNR
// is allowed to be any finite value; a negative threshold simply makes
// the link easier to pass.
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

// DefaultWarnBandDB is the margin threshold below which a link is
// labelled "marginal" instead of "healthy".
const DefaultWarnBandDB = 3.0

// RequiredPowerForMinSNR returns the received power in dBm that would
// just meet the required SNR: Pr_required = N_dBm + SNR_min.
func RequiredPowerForMinSNR(noiseDBm, minSNRDB float64) float64 {
	return noiseDBm + minSNRDB
}
