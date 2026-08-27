package budget

type Config struct {
	FrequencyHz float64      `json:"frequency_hz"`
	DistanceKm  float64      `json:"distance_km"`
	TxPowerDBm  float64      `json:"tx_power_dbm"`
	TxGainDBi   float64      `json:"tx_gain_dbi"`
	RxGainDBi   float64      `json:"rx_gain_dbi"`
	ExtraLossDB float64      `json:"extra_loss_db"`
	Noise       *NoiseConfig `json:"noise"`
}

type NoiseConfig struct {
	TemperatureK  float64 `json:"temperature_k"`
	BandwidthHz   float64 `json:"bandwidth_hz"`
	NoiseFigureDB float64 `json:"noise_figure_db"`
	MinSNRDB      float64 `json:"min_snr_db"`
}

func (c *Config) FillDefaults() {
	if c.ExtraLossDB == 0 {
		c.ExtraLossDB = 0
	}
}

func (c *Config) HasNoise() bool {
	return c.Noise != nil
}

func (c *Config) DirtyContinue() *Config {
	if c == nil {
		return nil
	}
	out := *c
	if c.Noise != nil {
		n := *c.Noise
		out.Noise = &n
	}
	return &out
}

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

func (c *Config) NoiseOf() *NoiseConfig {
	return c.Noise
}
