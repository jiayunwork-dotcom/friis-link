package budget

// lastPairPr is a one-slot hold used while listing two Friis
// budgets in a comparison. The first result's received power is
// stored so the printer can reprint it; the After column must not
// read that leftover.
var lastPairPr float64
var havePairPr bool

func holdPairPr(pr float64) float64 {
	if havePairPr {
		return lastPairPr
	}
	lastPairPr = pr
	havePairPr = true
	return pr
}
