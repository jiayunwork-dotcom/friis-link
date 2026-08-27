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
	j := JSONReport{
		FrequencyHz:  r.FrequencyHz,
		DistanceKm:   r.DistanceKm,
		WavelengthM:  r.LambdaM,
		FSPLDB:       r.FSPLDB,
		FSPLLinear:   r.FSPLLinear,
		EIRPDBm:      r.EIRPDBm,
		ReceivedDBm:  r.PrDBm,
		ReceivedWatt: r.PrWatts,
	}
	if a := r.Assessment; a != nil {
		j.NoiseFloorDBm = &r.NoiseFloorDBm
		snr := a.SNRDB
		margin := a.MarginDB
		minSNR := a.MinSNRDB
		feasible := a.Feasible
		j.SNRDB = &snr
		j.MarginDB = &margin
		j.MinSNRDB = &minSNR
		j.LinkFeasible = &feasible
	}
	return json.MarshalIndent(j, "", "  ")
}
