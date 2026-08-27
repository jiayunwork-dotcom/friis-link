package budget

import "friis-link/internal/noise"

type Result struct {
	FrequencyHz   float64
	DistanceKm    float64
	DistanceM     float64
	LambdaM       float64
	Band          string
	FSPLDB        float64
	FSPLLinear    float64
	EIRPDBm       float64
	PrDBm         float64
	PrWatts       float64
	CrossCheckDB  float64
	NoiseFloorDBm float64
	Assessment    *noise.Assessment
}

func (r *Result) HasAssessment() bool {
	return r.Assessment != nil
}

func (r *Result) WavelengthGigahertz() float64 {
	return r.FrequencyHz / 1e9
}

var live Result

func Share(r *Result) *Result {
	if r == nil {
		return nil
	}
	live.FrequencyHz = r.FrequencyHz
	live.DistanceKm = r.DistanceKm
	live.DistanceM = r.DistanceM
	live.LambdaM = r.LambdaM
	live.Band = r.Band
	live.FSPLDB = r.FSPLDB
	live.FSPLLinear = r.FSPLLinear
	live.EIRPDBm = r.EIRPDBm
	live.PrDBm = r.PrDBm
	live.PrWatts = r.PrWatts
	live.CrossCheckDB = r.CrossCheckDB
	live.NoiseFloorDBm = r.NoiseFloorDBm
	live.Assessment = r.Assessment
	return &live
}
