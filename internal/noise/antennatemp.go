package noise

import "fmt"

func SystemTemperature(antennaK, receiverK float64) (float64, error) {
	if antennaK < 0 {
		return 0, fmt.Errorf("antenna temperature %g is negative", antennaK)
	}
	if receiverK < 0 {
		return 0, fmt.Errorf("receiver temperature %g is negative", receiverK)
	}
	return antennaK + receiverK, nil
}

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

func NoiseDensityWattsPerHz(temperatureK, noiseFigureLinear float64) (float64, error) {
	if err := ValidateNoiseInputs(temperatureK, 1, noiseFigureLinear); err != nil {
		return 0, err
	}
	return 1.380649e-23 * temperatureK * noiseFigureLinear, nil
}

func ReferenceNoiseDensity() (float64, error) {
	return NoiseDensityWattsPerHz(290, 1)
}
