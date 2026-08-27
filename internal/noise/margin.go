package noise

func Margin(snrDB, minSNRDB float64) float64 {
	return snrDB - minSNRDB
}

func Feasible(marginDB float64) bool {
	return marginDB >= 0
}

type MarginState int

const (
	MarginInsufficient MarginState = iota
	MarginMarginal
	MarginHealthy
)

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

func ClassifyMargin(marginDB, warnBandDB float64) MarginState {
	if marginDB < 0 {
		return MarginInsufficient
	}
	if marginDB < warnBandDB {
		return MarginMarginal
	}
	return MarginHealthy
}
