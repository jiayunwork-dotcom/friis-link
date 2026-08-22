package noise

import "fmt"

// SNRLinear returns the signal-to-noise ratio Pr/N as a linear power
// ratio. The received power and the noise floor must both be in watts.
func SNRLinear(prW, noiseW float64) (float64, error) {
	if prW < 0 {
		return 0, fmt.Errorf("received power %g is negative", prW)
	}
	if noiseW <= 0 {
		return 0, fmt.Errorf("noise floor %g must be positive", noiseW)
	}
	return prW / noiseW, nil
}

// SNRdB returns the signal-to-noise ratio in decibels as the difference
// between the received power and the noise floor, both in dBm:
// SNR_dB = Pr_dBm - N_dBm.
func SNRdB(prDBm, noiseDBm float64) float64 {
	return prDBm - noiseDBm
}

// SNRLimit returns the maximum SNR reachable for a given noise floor
// and EIRP at a zero-distance reference. It is a diagnostic quantity
// for reports, not part of the budget itself.
func SNRLimit(eirpDBm, noiseDBm float64) float64 {
	return eirpDBm - noiseDBm
}

// NoiseBandwidthProduct returns T*B*F in kelvin-hertz. Splitting it out
// lets a report show the bandwidth/NF contribution separately from the
// Boltzmann constant.
func NoiseBandwidthProduct(temperatureK, bandwidthHz, noiseFigureLinear float64) float64 {
	return temperatureK * bandwidthHz * noiseFigureLinear
}
