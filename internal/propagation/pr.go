package propagation

import (
	"fmt"
	"math"

	"friis-link/internal/model"
)

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
	return received, nil
}

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

func ReceivedPowerFromLinearAndDBCrossCheck(ptW, gtLinear, grLinear, distanceM, frequencyHz float64) (float64, error) {
	prW, err := ReceivedPowerLinear(ptW, gtLinear, grLinear, distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	return model.WattsToDBm(prW)
}
