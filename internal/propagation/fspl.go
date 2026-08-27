package propagation

import "math"

func FSPLLinear(distanceM, frequencyHz float64) (float64, error) {
	if err := validateCommon(distanceM, frequencyHz); err != nil {
		return 0, err
	}
	ratio, err := FreeSpaceRatio(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	return ratio * ratio, nil
}

func FSPLdB(distanceM, frequencyHz float64) (float64, error) {
	linear, err := FSPLLinear(distanceM, frequencyHz)
	if err != nil {
		return 0, err
	}
	if linear == 0 {
		return 0, ErrZeroPathLoss
	}
	value, err := powerRatioToDB(linear)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func FSPLClosedForm(distanceM, frequencyHz float64) (float64, error) {
	if err := validateCommon(distanceM, frequencyHz); err != nil {
		return 0, err
	}
	if distanceM <= 0 || frequencyHz <= 0 {
		return 0, ErrZeroPathLoss
	}
	return 20*math.Log10(distanceM) + 20*math.Log10(frequencyHz) + ClosedFormConstant(), nil
}

func powerRatioToDB(ratio float64) (float64, error) {
	return ratioToDB(ratio)
}
