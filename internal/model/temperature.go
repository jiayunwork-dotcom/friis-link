package model

import "fmt"

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func KelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

func FormatTemperature(kelvin float64) string {
	return fmt.Sprintf("%.4g K", kelvin)
}
