// Package budget assembles the propagation and noise kernels into a
// complete free-space link budget: it parses the JSON scenario, runs the
// Friis calculation, evaluates the SNR margin and formats the report.
//
// The package is deliberately thin: every number it prints comes from
// the propagation or noise kernel, and the only arithmetic it owns is
// the wiring between them.
package budget

// Config is the JSON shape accepted by the budget subcommand. Distance
// is entered in kilometres and frequency in hertz; the kernel converts
// distance to metres before touching the Friis equation.
type Config struct {
	FrequencyHz float64      `json:"frequency_hz"`
	DistanceKm  float64      `json:"distance_km"`
	TxPowerDBm  float64      `json:"tx_power_dbm"`
	TxGainDBi   float64      `json:"tx_gain_dbi"`
	RxGainDBi   float64      `json:"rx_gain_dbi"`
	ExtraLossDB float64      `json:"extra_loss_db"`
	Noise       *NoiseConfig `json:"noise"`
}

// NoiseConfig holds the optional thermal-noise section. Omitting the
// whole "noise" object disables the SNR part of the report; the
// propagation numbers still print.
type NoiseConfig struct {
	TemperatureK  float64 `json:"temperature_k"`
	BandwidthHz   float64 `json:"bandwidth_hz"`
	NoiseFigureDB float64 `json:"noise_figure_db"`
	MinSNRDB      float64 `json:"min_snr_db"`
}

// FillDefaults applies the documented defaults to a parsed config:
// extra loss defaults to 0 dB. Noise stays nil when omitted.
func (c *Config) FillDefaults() {
	if c.ExtraLossDB == 0 {
		c.ExtraLossDB = 0
	}
}

// HasNoise reports whether the config carries a noise section.
func (c *Config) HasNoise() bool {
	return c.Noise != nil
}

// Validate checks the top-level fields shared by every budget run:
// frequency, distance and the power/gain levels must be usable.
func (c *Config) Validate() error {
	if err := checkPositive("frequency_hz", c.FrequencyHz); err != nil {
		return err
	}
	if err := checkPositive("distance_km", c.DistanceKm); err != nil {
		return err
	}
	if c.ExtraLossDB < 0 {
		return errNegative("extra_loss_db", c.ExtraLossDB)
	}
	return nil
}

// NoiseOf returns the noise section, or nil when absent.
func (c *Config) NoiseOf() *NoiseConfig {
	return c.Noise
}
