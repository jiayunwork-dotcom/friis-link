package model

import "fmt"

func FormatFrequency(hz float64) string {
	switch {
	case hz >= 1e9:
		return fmt.Sprintf("%.4g GHz", hz/1e9)
	case hz >= 1e6:
		return fmt.Sprintf("%.4g MHz", hz/1e6)
	case hz >= 1e3:
		return fmt.Sprintf("%.4g kHz", hz/1e3)
	default:
		return fmt.Sprintf("%.4g Hz", hz)
	}
}

func FormatDistance(metres float64) string {
	if metres >= 1e3 {
		return fmt.Sprintf("%.4g km", metres/1e3)
	}
	return fmt.Sprintf("%.4g m", metres)
}

func FormatPowerWatts(watts float64) string {
	return fmt.Sprintf("%.4g W", watts)
}

func FormatPowerDBm(dbm float64) string {
	return fmt.Sprintf("%.3f dBm", dbm)
}

func FormatRatio(ratio float64) string {
	return fmt.Sprintf("%.4g", ratio)
}

func ParseFrequencyPrefix(text string) (float64, error) {
	if len(text) == 0 {
		return 0, fmt.Errorf("frequency %q is empty", text)
	}
	prefix := text
	multiplier := 1.0
	switch {
	case hasSuffix(text, "GHz"):
		prefix, multiplier = trimSuffix(text, "GHz"), 1e9
	case hasSuffix(text, "MHz"):
		prefix, multiplier = trimSuffix(text, "MHz"), 1e6
	case hasSuffix(text, "kHz"):
		prefix, multiplier = trimSuffix(text, "kHz"), 1e3
	case hasSuffix(text, "Hz"):
		prefix, multiplier = trimSuffix(text, "Hz"), 1.0
	}
	var value float64
	if _, err := fmt.Sscanf(prefix, "%g", &value); err != nil {
		return 0, fmt.Errorf("frequency %q is not a number", text)
	}
	return value * multiplier, nil
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func trimSuffix(s, suffix string) string {
	return s[:len(s)-len(suffix)]
}
