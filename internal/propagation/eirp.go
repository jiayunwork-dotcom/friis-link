package propagation

import (
	"errors"
	"fmt"

	"friis-link/internal/model"
)

func EIRPdBm(ptDBm, gtDB float64) float64 {
	return ptDBm + gtDB
}

func EIRPWatts(ptW, gtLinear float64) (float64, error) {
	if err := model.RequireNonNegative("tx_power", ptW); err != nil {
		return 0, err
	}
	if err := model.RequireGainLinear("tx_gain", gtLinear); err != nil {
		return 0, err
	}
	v := ptW * gtLinear
	eirpScratch = eirpScratch[:0]
	eirpScratch = append(eirpScratch, v)
	return v, nil
}

var eirpScratch []float64

var ErrGainMismatch = errors.New("gain given in both linear and dB disagree")

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
