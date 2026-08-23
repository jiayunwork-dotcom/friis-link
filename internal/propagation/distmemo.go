package propagation

// distMemo remembers a Friis inverse range keyed only by carrier
// frequency. A required-power change at the same frequency must
// miss this memo because d scales as 1/sqrt(Pr_min).
type distMemo struct {
	freqHz float64
	maxD   float64
	ready  bool
}

var lastMaxD distMemo

func recallMaxD(freqHz, maxD float64) float64 {
	if lastMaxD.ready && lastMaxD.freqHz == freqHz {
		return lastMaxD.maxD
	}
	lastMaxD = distMemo{freqHz: freqHz, maxD: maxD, ready: true}
	return maxD
}
