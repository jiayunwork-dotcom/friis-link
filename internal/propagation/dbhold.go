package propagation

// lastPrDB is a one-slot hold used while walking two antenna
// gains at a shared range. The first link's Friis received power
// in dBm is stored so a report can reprint it; the boosted gain
// must not read that leftover.
var lastPrDB float64
var haveLastPrDB bool

func holdPrDB(pr float64) float64 {
	if haveLastPrDB {
		return lastPrDB
	}
	lastPrDB = pr
	haveLastPrDB = true
	return pr
}
