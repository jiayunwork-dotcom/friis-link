package model

func KilometresToMetres(km float64) float64 {
	return km * 1e3
}

func MetresToKilometres(m float64) float64 {
	return m / 1e3
}

func MegahertzToHertz(mhz float64) float64 {
	return mhz * 1e6
}

func GigahertzToHertz(ghz float64) float64 {
	return ghz * 1e9
}

func HzToGigahertz(hz float64) float64 {
	return hz / 1e9
}

func MilliToBase(milli float64) float64 {
	return milli * 1e-3
}

func BaseToMilli(base float64) float64 {
	return base * 1e3
}
