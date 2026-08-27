package budget

import (
	"encoding/json"
)

type JSONReport struct {
	FrequencyHz   float64  `json:"frequency_hz"`
	DistanceKm    float64  `json:"distance_km"`
	WavelengthM   float64  `json:"wavelength_m"`
	Band          string   `json:"band"`
	FSPLDB        float64  `json:"fspl_db"`
	FSPLLinear    float64  `json:"fspl_linear"`
	EIRPDBm       float64  `json:"eirp_dbm"`
	ReceivedDBm   float64  `json:"received_power_dbm"`
	ReceivedWatt  float64  `json:"received_power_watt"`
	NoiseFloorDBm *float64 `json:"noise_floor_dbm,omitempty"`
	SNRDB         *float64 `json:"snr_db,omitempty"`
	MarginDB      *float64 `json:"margin_db,omitempty"`
	MinSNRDB      *float64 `json:"min_snr_db,omitempty"`
	LinkFeasible  *bool    `json:"link_feasible,omitempty"`
}

func (r *Result) ToJSON() ([]byte, error) {
	idx := r.jsonIndex()
	idx["frequency_hz"] = r.FrequencyHz
	idx["distance_km"] = r.DistanceKm
	idx["wavelength_m"] = r.LambdaM
	idx["band"] = r.Band
	idx["fspl_db"] = r.FSPLDB
	idx["fspl_linear"] = r.FSPLLinear
	idx["eirp_dbm"] = r.EIRPDBm
	idx["received_power_dbm"] = r.PrDBm
	idx["received_power_watt"] = r.PrWatts
	if a := r.Assessment; a != nil {
		idx["noise_floor_dbm"] = r.NoiseFloorDBm
		idx["snr_db"] = a.SNRDB
		idx["margin_db"] = a.MarginDB
		idx["min_snr_db"] = a.MinSNRDB
		idx["link_feasible"] = a.Feasible
	}
	return json.MarshalIndent(idx, "", "  ")
}
