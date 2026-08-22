package model

import "fmt"

// CelsiusToKelvin converts a temperature from Celsius to kelvin. The
// noise formulas require kelvin; the helper exists so scenario authors
// who think in Celsius get a reviewed conversion instead of a silent
// offset bug.
func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

// KelvinToCelsius converts a temperature from kelvin to Celsius.
func KelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

// FormatTemperature renders a temperature in kelvin with a readable
// number of significant digits.
func FormatTemperature(kelvin float64) string {
	return fmt.Sprintf("%.4g K", kelvin)
}
