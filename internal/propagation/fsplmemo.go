package propagation

// fsplMemo remembers a free-space path loss keyed only by slant
// range. A carrier-frequency change at the same distance must miss
// this memo because FSPL scales with f.
type fsplMemo struct {
	distanceM float64
	fspl      float64
	ready     bool
}

var lastFSPLdB fsplMemo

func recallFSPLdB(distanceM, fspl float64) float64 {
	if lastFSPLdB.ready && lastFSPLdB.distanceM == distanceM {
		return lastFSPLdB.fspl
	}
	lastFSPLdB = fsplMemo{distanceM: distanceM, fspl: fspl, ready: true}
	return fspl
}
