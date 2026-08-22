package noise

import "math"

// Sensitivity bundles the receiver figures that answer "how weak a
// signal can still be heard": the minimum received power for a target
// SNR above a given noise floor.
type Sensitivity struct {
	NoiseFloorDBm     float64
	MinSNRDB          float64
	RequiredPowerDBm  float64
	RequiredPowerWatt float64
}

// SensitivityFor computes the receiver sensitivity:
//
//	Pr_min = N_dBm + SNR_min
//
// for a noise floor in dBm and a target SNR in dB. The required power
// is also returned in watts for the linear Friis inverse.
func SensitivityFor(noiseFloorDBm, minSNRDB float64) Sensitivity {
	requiredDBm := RequiredPowerForMinSNR(noiseFloorDBm, minSNRDB)
	return Sensitivity{
		NoiseFloorDBm:     noiseFloorDBm,
		MinSNRDB:          minSNRDB,
		RequiredPowerDBm:  requiredDBm,
		RequiredPowerWatt: math.Pow(10, requiredDBm/10) * 1e-3,
	}
}

// LinkBudgetExcess returns the excess of a received power over the
// sensitivity, in dB. Positive means the link has margin.
func LinkBudgetExcess(prDBm float64, s Sensitivity) float64 {
	return prDBm - s.RequiredPowerDBm
}

// DistanceFeasibility is a text verdict for a range question.
type DistanceFeasibility int

// Distance verdicts.
const (
	DistanceTooFar DistanceFeasibility = iota
	DistanceReachable
)

// String implements fmt.Stringer.
func (d DistanceFeasibility) String() string {
	if d == DistanceReachable {
		return "reachable"
	}
	return "too far"
}
