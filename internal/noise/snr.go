package noise

import "fmt"

func SNRLinear(prW, noiseW float64) (float64, error) {
	if prW < 0 {
		return 0, fmt.Errorf("received power %g is negative", prW)
	}
	if noiseW <= 0 {
		return 0, fmt.Errorf("noise floor %g must be positive", noiseW)
	}
	return prW / noiseW, nil
}

func SNRdB(prDBm, noiseDBm float64) float64 {
	return prDBm - noiseDBm
}

func SNRLimit(eirpDBm, noiseDBm float64) float64 {
	return eirpDBm - noiseDBm
}

func NoiseBandwidthProduct(temperatureK, bandwidthHz, noiseFigureLinear float64) float64 {
	return temperatureK * bandwidthHz * noiseFigureLinear
}
