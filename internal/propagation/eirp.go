package propagation

import (
	"errors"
	"fmt"

	"friis-link/internal/model"
)

// EIRPdBm returns the effective isotropic radiated power in dBm,
// EIRP = Pt + Gt. The identity EIRP_dBm = Pt_dBm + Gt_dB is pinned by a
// test so the sum can never be reordered into a product by accident.
func EIRPdBm(ptDBm, gtDB float64) float64 {
	return ptDBm + gtDB
}

// EIRPWatts returns the effective isotropic radiated power in watts:
// Pt * Gt with both quantities linear.
func EIRPWatts(ptW, gtLinear float64) (float64, error) {
	if err := model.RequireNonNegative("tx_power", ptW); err != nil {
		return 0, err
	}
	if err := model.RequireGainLinear("tx_gain", gtLinear); err != nil {
		return 0, err
	}
	return ptW * gtLinear, nil
}

// ErrGainMismatch is returned when a linear gain and its decibel
// counterpart are inconsistent for the same antenna.
var ErrGainMismatch = errors.New("gain given in both linear and dB disagree")

// CheckGainConsistency verifies that a linear gain converts to the
// declared decibel gain. Used by the budget kernel when a config carries
// both representations, so a hand-edited JSON cannot silently disagree.
func CheckGainConsistency(gainLinear, gainDB float64) error {
	converted, err := model.GainLinearToDB(gainLinear)
	if err != nil {
		return err
	}
	const tolerance = 1e-9
	if abs(converted-gainDB) > tolerance {
		return fmt.Errorf("%w: linear %g -> %.6f dB, declared %.6f dB",
			ErrGainMismatch, gainLinear, converted, gainDB)
	}
	return nil
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
