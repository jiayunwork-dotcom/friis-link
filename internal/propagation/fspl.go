package propagation

import "math"

// FSPLLinear returns the free-space path loss as a linear power ratio
// (4*pi*d/lambda)^2. Doubling the distance quadruples the ratio, which
// shows up as -6 dB in the decibel domain.
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

// FSPLdB returns the free-space path loss in decibels,
// 20*log10(4*pi*d/lambda). The decibel form is derived from the same
// ratio as FSPLLinear so the two can never disagree by a 10log/20log
// factor.
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

// FSPLClosedForm evaluates 20*log10(d) + 20*log10(f) + C with d in
// metres and f in hertz. It must agree with FSPLdB; the test
// TestFSPLMatchesClosedForm pins that agreement so a refactor cannot
// silently split the two implementations apart.
func FSPLClosedForm(distanceM, frequencyHz float64) (float64, error) {
	if err := validateCommon(distanceM, frequencyHz); err != nil {
		return 0, err
	}
	if distanceM <= 0 || frequencyHz <= 0 {
		return 0, ErrZeroPathLoss
	}
	return 20*math.Log10(distanceM) + 20*math.Log10(frequencyHz) + ClosedFormConstant(), nil
}

// powerRatioToDB converts a linear power ratio to decibels using the
// single shared conversion owned by the model package.
func powerRatioToDB(ratio float64) (float64, error) {
	return ratioToDB(ratio)
}
