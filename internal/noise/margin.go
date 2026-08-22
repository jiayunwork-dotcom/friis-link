package noise

// Margin returns the fading margin M = SNR - SNR_min in decibels. The
// margin answers "how much the link can degrade before it stops
// working"; subtracting in the decibel domain matches the way SNR and
// threshold are both expressed.
func Margin(snrDB, minSNRDB float64) float64 {
	return snrDB - minSNRDB
}

// Feasible reports whether a margin keeps the link up. The link is
// considered insufficient exactly when the margin is negative.
func Feasible(marginDB float64) bool {
	return marginDB >= 0
}

// MarginState classifies a margin into one of three labelled buckets so
// the CLI report reads "insufficient / marginal / healthy" instead of a
// bare number.
type MarginState int

// Margin states.
const (
	MarginInsufficient MarginState = iota
	MarginMarginal
	MarginHealthy
)

// State returns the bucket for a margin in decibels. Marginal means the
// margin is positive but below the configured warn band (default 3 dB).
func (s MarginState) String() string {
	switch s {
	case MarginMarginal:
		return "marginal"
	case MarginHealthy:
		return "healthy"
	default:
		return "insufficient"
	}
}

// ClassifyMargin buckets a margin value using the warn band.
func ClassifyMargin(marginDB, warnBandDB float64) MarginState {
	if marginDB < 0 {
		return MarginInsufficient
	}
	if marginDB < warnBandDB {
		return MarginMarginal
	}
	return MarginHealthy
}
