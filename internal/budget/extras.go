package budget

import (
	"fmt"

	"friis-link/internal/gas"
	"friis-link/internal/pointing"
	"friis-link/internal/rain"
	"friis-link/internal/slant"
)

func Clone(cfg *Config) *Config {
	if cfg == nil {
		return nil
	}
	out := *cfg
	if cfg.Noise != nil {
		n := *cfg.Noise
		out.Noise = &n
	}
	return &out
}

func WithExtras(cfg *Config, extras ...float64) (*Config, error) {
	if cfg == nil {
		return nil, errNilConfig
	}
	sum, err := pointing.CombineOnce(extras...)
	if err != nil {
		return nil, err
	}
	out := Clone(cfg)
	out.ExtraLossDB = cfg.ExtraLossDB + sum
	return out, nil
}

func ComputeWithExtrasOnce(cfg *Config, extras ...float64) (*Result, error) {
	aug, err := WithExtras(cfg, extras...)
	if err != nil {
		return nil, err
	}
	return Compute(aug)
}

func RainOnSlant(cfg *Config, k, alpha, rateMMH, elevDeg, reduction float64) (float64, error) {
	if cfg == nil {
		return 0, errNilConfig
	}
	return slant.RainPath(k, alpha, rateMMH, cfg.DistanceKm, elevDeg, reduction)
}

func GasOnPath(cfg *Config, a gas.Air) (float64, error) {
	if cfg == nil {
		return 0, errNilConfig
	}
	return gas.Path(cfg.FrequencyHz/1e9, cfg.DistanceKm, a)
}

func ComputeWithRainGasOnce(cfg *Config, k, alpha, rateMMH, elevDeg, reduction float64, a gas.Air) (*Result, error) {
	rainDB, err := RainOnSlant(cfg, k, alpha, rateMMH, elevDeg, reduction)
	if err != nil {
		return nil, err
	}
	gasDB, err := GasOnPath(cfg, a)
	if err != nil {
		return nil, err
	}
	return ComputeWithExtrasOnce(cfg, rainDB, gasDB)
}

func DoubleApplyDetected(cfg *Config, extra float64) error {
	once, err := ComputeWithExtrasOnce(cfg, extra)
	if err != nil {
		return err
	}
	twice, err := WithExtras(cfg, extra)
	if err != nil {
		return err
	}
	twice, err = WithExtras(twice, extra)
	if err != nil {
		return err
	}
	dup, err := Compute(twice)
	if err != nil {
		return err
	}
	if dup.PrDBm >= once.PrDBm-0.01 {
		return fmt.Errorf("budget: applying extra twice should drop Pr: once %g twice %g", once.PrDBm, dup.PrDBm)
	}
	delta := once.PrDBm - dup.PrDBm
	if delta < extra-1e-6 {
		return fmt.Errorf("budget: double apply should cost another %g dB, got %g", extra, delta)
	}
	return nil
}

func HigherElevationCutsRain(cfg *Config, k, alpha, rateMMH, lowDeg, highDeg, reduction float64) error {
	if cfg == nil {
		return errNilConfig
	}
	return slant.HigherElevationShortensRain(k, alpha, rateMMH, cfg.DistanceKm, lowDeg, highDeg, reduction)
}

func SlantRainExceedsGroundSlice(cfg *Config, k, alpha, rateMMH, elevDeg float64) error {
	if cfg == nil {
		return errNilConfig
	}
	ground, err := rain.Path(k, alpha, rateMMH, cfg.DistanceKm)
	if err != nil {
		return err
	}
	sl, err := RainOnSlant(cfg, k, alpha, rateMMH, elevDeg, 1)
	if err != nil {
		return err
	}
	if elevDeg < 90 && sl <= ground {
		return fmt.Errorf("budget: slant rain %g should exceed the ground-length slice %g", sl, ground)
	}
	return nil
}
