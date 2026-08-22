package budget

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// ParseConfig decodes a scenario from raw JSON. Unknown fields are
// rejected so a typo in the input cannot silently change the link: an
// operator who writes "frequency_Hz" instead of "frequency_hz" gets an
// error instead of a budget computed with a zero frequency.
func ParseConfig(data []byte) (*Config, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.FillDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if n := cfg.Noise; n != nil {
		if err := n.Validate(); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

// LoadConfig reads and parses a scenario file from disk.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	cfg, err := ParseConfig(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

// Validate checks the noise section. A zero bandwidth is rejected so
// the SNR cannot divide by zero; a non-positive temperature or noise
// figure is equally meaningless.
func (n *NoiseConfig) Validate() error {
	if err := checkPositive("temperature_k", n.TemperatureK); err != nil {
		return err
	}
	if err := checkPositive("bandwidth_hz", n.BandwidthHz); err != nil {
		return err
	}
	if err := checkPositive("noise_figure_db", n.NoiseFigureDB); err != nil {
		return err
	}
	return nil
}

// checkPositive mirrors the kernel validation for a named quantity.
func checkPositive(name string, v float64) error {
	if v <= 0 {
		return fmt.Errorf("%s: %g must be strictly positive", name, v)
	}
	return nil
}

func errNegative(name string, v float64) error {
	return fmt.Errorf("%s: %g must be non-negative", name, v)
}
