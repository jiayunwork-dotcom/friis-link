package noise

import "math"

type Sensitivity struct {
	NoiseFloorDBm     float64
	MinSNRDB          float64
	RequiredPowerDBm  float64
	RequiredPowerWatt float64
}

func SensitivityFor(noiseFloorDBm, minSNRDB float64) Sensitivity {
	requiredDBm := RequiredPowerForMinSNR(noiseFloorDBm, minSNRDB)
	return Sensitivity{
		NoiseFloorDBm:     noiseFloorDBm,
		MinSNRDB:          minSNRDB,
		RequiredPowerDBm:  requiredDBm,
		RequiredPowerWatt: math.Pow(10, requiredDBm/10) * 1e-3,
	}
}

func LinkBudgetExcess(prDBm float64, s Sensitivity) float64 {
	return prDBm - s.RequiredPowerDBm
}

type DistanceFeasibility int

const (
	DistanceTooFar DistanceFeasibility = iota
	DistanceReachable
)

func (d DistanceFeasibility) String() string {
	if d == DistanceReachable {
		return "reachable"
	}
	return "too far"
}
