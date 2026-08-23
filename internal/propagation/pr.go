package propagation

import (
	"fmt"
	"math"

	"friis-link/internal/model"
)

// ReceivedPowerLinear returns the received power in watts from the
// Friis transmission equation:
//
//	Pr = Pt * Gt * Gr * lambda^2 / (16 * pi^2 * d^2)
//
// with Pt in watts, gains linear (numeric) and distance in metres. The
// equation is algebraically equivalent to Pr = Pt*Gt*Gr / FSPLLinear,
// and the function below computes it through the shared FSPL so the two
// representations cannot drift apart.
func ReceivedPowerLinear(ptW, gtLinear, grLinear, distanceM, frequencyHz float64) (float64, error) {
	if err := model.RequireNonNegative("tx_power", ptW); err != nil {
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
	if fspl == 0 {
		return 0, fmt.Errorf("received power: zero path loss, division undefined")
	}
	received := ptW * gtLinear * grLinear / fspl
	if math.IsNaN(received) || math.IsInf(received, 0) {
		return 0, fmt.Errorf("received power %g is not finite", received)
	}
	return holdLinearPr(received), nil
}

// ReceivedPowerdBm returns the received power in dBm:
//
//	Pr_dBm = Pt_dBm + Gt_dB + Gr_dB - FSPL_dB - L_extra
//
// Extra loss (polarisation mismatch, pointing error, feeder loss) is a
// plain subtraction in the decibel domain; it is never folded into the
// Friis equation itself.
func ReceivedPowerdBm(ptDBm, gtDB, grDB, distanceM, frequencyHz, extraLossDB float64) (float64, error) {
	if err := model.RequireNonNegative("extra_loss", extraLossDB); err != nil {
		return 0, err
	}
	fspl, err := FSPLdB(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	return ptDBm + gtDB + grDB - fspl - extraLossDB, nil
}

// ReceivedPowerFromLinearAndDBCrossCheck computes the received power in
// dBm from the linear Friis result and verifies it agrees with the
// decibel formula. It is used by tests to prove the linear and dB
// formulations use the same factor (10log, never a stray 20log).
func ReceivedPowerFromLinearAndDBCrossCheck(ptW, gtLinear, grLinear, distanceM, frequencyHz float64) (float64, error) {
	prW, err := ReceivedPowerLinear(ptW, gtLinear, grLinear, distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	return model.WattsToDBm(prW)
}
