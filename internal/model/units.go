package model

// Unit conversions used by the budget input. The CLI accepts distances
// in kilometres and frequencies in hertz (GHz/MHz are resolved by the
// caller before reaching the kernel); every internal formula operates on
// metres and hertz.

// KilometresToMetres scales a distance from km to m.
func KilometresToMetres(km float64) float64 {
	return km * 1e3
}

// MetresToKilometres scales a distance from m to km.
func MetresToKilometres(m float64) float64 {
	return m / 1e3
}

// MegahertzToHertz scales a frequency from MHz to Hz.
func MegahertzToHertz(mhz float64) float64 {
	return mhz * 1e6
}

// GigahertzToHertz scales a frequency from GHz to Hz.
func GigahertzToHertz(ghz float64) float64 {
	return ghz * 1e9
}

// HzToGigahertz scales a frequency from Hz to GHz.
func HzToGigahertz(hz float64) float64 {
	return hz / 1e9
}

// MilliToBase converts a milli-scaled quantity (milliwatt, milliwatt
// per hertz) to its base unit.
func MilliToBase(milli float64) float64 {
	return milli * 1e-3
}

// BaseToMilli converts a base-unit quantity to the milli scale.
func BaseToMilli(base float64) float64 {
	return base * 1e3
}
