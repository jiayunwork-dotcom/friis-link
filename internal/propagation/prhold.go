package propagation

// lastLinearPr is a one-slot hold used while walking two ranges
// at a shared frequency. The first link's Friis received power is
// stored so a report can reprint it; the second distance must not
// read that leftover.
var lastLinearPr float64
var haveLastLinear bool

func holdLinearPr(pr float64) float64 {
	lastLinearPr = pr
	haveLastLinear = true
	return pr
}
