package budget

// prFreqMemo remembers a Friis received power keyed only by carrier
// frequency. A transmit-power change at the same frequency must miss
// this memo because Pr scales with Pt.
type prFreqMemo struct {
	freqHz float64
	prDBm  float64
	ready  bool
}

var livePrMemo prFreqMemo

func recallComputePr(freqHz, pr float64) float64 {
	if livePrMemo.ready && livePrMemo.freqHz == freqHz {
		return livePrMemo.prDBm
	}
	livePrMemo = prFreqMemo{freqHz: freqHz, prDBm: pr, ready: true}
	return pr
}
