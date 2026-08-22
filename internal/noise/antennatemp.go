package noise

import "fmt"

// SystemTemperature combines the antenna noise temperature and the
// receiver noise temperature into the system temperature used by the
// kTBF formula:
//
//	T_sys = T_antenna + T_receiver
func SystemTemperature(antennaK, receiverK float64) (float64, error) {
	if antennaK < 0 {
		return 0, fmt.Errorf("antenna temperature %g is negative", antennaK)
	}
	if receiverK < 0 {
		return 0, fmt.Errorf("receiver temperature %g is negative", receiverK)
	}
	return antennaK + receiverK, nil
}

// TemperatureFromNoiseFloor recovers the temperature that a given noise
// power implies for a bandwidth and noise figure: T = N/(k*B*F). Useful
// for cross-checking a measured floor against the model.
func TemperatureFromNoiseFloor(noiseW, bandwidthHz, noiseFigureLinear float64) (float64, error) {
	if err := ValidateNoiseInputs(290, bandwidthHz, noiseFigureLinear); err != nil {
		return 0, err
	}
	if noiseW < 0 {
		return 0, fmt.Errorf("noise power %g is negative", noiseW)
	}
	denominator := 1.380649e-23 * bandwidthHz * noiseFigureLinear
	if denominator == 0 {
		return 0, fmt.Errorf("temperature: zero k*B*F denominator")
	}
	return noiseW / denominator, nil
}

// NoiseDensityWattsPerHz returns the one-sided noise power spectral
// density N0 = k*T*F in watts per hertz, independent of bandwidth.
func NoiseDensityWattsPerHz(temperatureK, noiseFigureLinear float64) (float64, error) {
	if err := ValidateNoiseInputs(temperatureK, 1, noiseFigureLinear); err != nil {
		return 0, err
	}
	return 1.380649e-23 * temperatureK * noiseFigureLinear, nil
}

// ReferenceNoiseDensity returns N0 at the 290 K reference temperature
// with a unit noise figure: about -174 dBm/Hz.
func ReferenceNoiseDensity() (float64, error) {
	return NoiseDensityWattsPerHz(290, 1)
}
