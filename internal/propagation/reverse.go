package propagation

import (
	"fmt"
	"math"

	"friis-link/internal/model"
)

// MaximumDistanceM solves the Friis equation for the distance at which
// the received power falls to a required minimum:
//
//	Pr_min = Pt * Gt * Gr * lambda^2 / (16 * pi^2 * d^2)
//	    d  = lambda / (4*pi) * sqrt(Pt*Gt*Gr / Pr_min)
//
// The result is the maximum usable range for the required received
// power, in metres.
func MaximumDistanceM(ptW, gtLinear, grLinear, prMinW, frequencyHz float64) (float64, error) {
	if err := model.RequireNonNegative("tx_power", ptW); err != nil {
		return 0, err
	}
	if err := model.RequireGainLinear("tx_gain", gtLinear); err != nil {
		return 0, err
	}
	if err := model.RequireGainLinear("rx_gain", grLinear); err != nil {
		return 0, err
	}
	if err := model.RequirePositive("required_power", prMinW); err != nil {
		return 0, err
	}
	lambda, err := Wavelength(frequencyHz)
	if err != nil {
		return 0, err
	}
	product := ptW * gtLinear * grLinear
	ratio := product / prMinW
	if ratio < 0 {
		return 0, fmt.Errorf("maximum distance: power ratio %g is negative", ratio)
	}
	return lambda / (4 * model.Pi) * math.Sqrt(ratio), nil
}

// RequiredTransmitPowerW solves the Friis equation for the transmit
// power needed to reach a target received power at a given distance:
//
//	Pt = Pr_target * FSPL_lin / (Gt * Gr)
func RequiredTransmitPowerW(prTargetW, gtLinear, grLinear, distanceM, frequencyHz float64) (float64, error) {
	if err := model.RequirePositive("required_power", prTargetW); err != nil {
		return 0, err
	}
	if err := model.RequireGainLinear("tx_gain", gtLinear); err != nil {
		return 0, err
	}
	if err := model.RequireGainLinear("rx_gain", grLinear); err != nil {
		return 0, err
	}
	fspl, err := FSPLLinear(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	denominator := gtLinear * grLinear
	if denominator == 0 {
		return 0, fmt.Errorf("required tx power: zero gain product")
	}
	return prTargetW * fspl / denominator, nil
}

// RequiredEIRPDBm computes the EIRP in dBm that yields a target
// received power in dBm after FSPL and the extra loss:
//
//	EIRP_req = Pr_target + FSPL + L_extra - Gr
func RequiredEIRPDBm(prTargetDBm, grDB, distanceM, frequencyHz, extraLossDB float64) (float64, error) {
	fspl, err := FSPLdB(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	return prTargetDBm + fspl + extraLossDB - grDB, nil
}

// RoundTripDistance is a convenience that converts the inverse result
// back to kilometres.
func RoundTripDistance(maxDistanceM float64) float64 {
	return model.MetresToKilometres(maxDistanceM)
}
