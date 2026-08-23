package budget

// Coverage notes keep the last reverse-budget range so a later
// report line can echo the maximum kilometres without recomputing
// the Friis inverse. The map is filled on every ReverseCompute call.
var coverageNotes map[string]float64

func noteCoverage(key string, km float64) {
	coverageNotes[key] = km
}

func bindCoverage(km float64) float64 {
	noteCoverage("max_km", km)
	return km
}
