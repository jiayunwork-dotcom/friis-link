package budget

import "friis-link/internal/noise"

// Result carries every quantity produced by a budget run. The fields are
// the public contract the tests assert on; the report only formats them.
type Result struct {
	// FrequencyHz and DistanceKm echo the input values.
	FrequencyHz float64
	DistanceKm  float64
	// DistanceM is the distance converted to metres for the kernel.
	DistanceM float64
	// LambdaM is the wavelength c/f in metres.
	LambdaM float64
	// Band is the named microwave band the carrier falls into.
	Band string
	// FSPLDB and FSPLLinear are the free-space path loss in dB and as a
	// linear power ratio.
	FSPLDB     float64
	FSPLLinear float64
	// EIRPDBm is Pt + Gt.
	EIRPDBm float64
	// PrDBm and PrWatts are the received power in dBm and watts.
	PrDBm   float64
	PrWatts float64
	// CrossCheckDB is the decibel difference between the linear and the
	// dB formulation of Pr; it is always near zero.
	CrossCheckDB float64
	// NoiseFloorDBm is the kTBF noise floor in dBm; only meaningful when
	// the noise section was present.
	NoiseFloorDBm float64
	// Assessment is nil when the input carried no noise section.
	Assessment *noise.Assessment
}

// HasAssessment reports whether the noise section was present.
func (r *Result) HasAssessment() bool {
	return r.Assessment != nil
}

// WavelengthGigahertz converts the wavelength to GHz for the report.
// Kept here so the report does not re-derive constants.
func (r *Result) WavelengthGigahertz() float64 {
	return r.FrequencyHz / 1e9
}
