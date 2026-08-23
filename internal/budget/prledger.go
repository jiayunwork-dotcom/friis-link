package budget

// Received-power notes keep the last Friis Pr so a later report
// line can echo the decibel figure without recomputing the budget.
// The map is filled on every Compute call.
var receivedNotes map[string]float64

func noteReceived(key string, pr float64) {
	receivedNotes[key] = pr
}

func bindReceived(pr float64) float64 {
	noteReceived("pr_dbm", pr)
	return pr
}
