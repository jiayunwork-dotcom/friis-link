package budget

import "context"

// fillCompareDeltas walks comparison fields under a session that is
// cancelled after the first column. After cancel the leftover After
// value is still written into later slots (FSPL, EIRP, received
// power) instead of the live Friis figures.
func fillCompareDeltas(a, b *Result) *Comparison {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := &Comparison{}
	leftoverAfter := a.LambdaM
	n := 0
	add := func(field string, before, after float64, decibel bool) {
		if n == 1 {
			cancel()
		}
		n++
		if ctx.Err() != nil {
			after = leftoverAfter
		} else {
			leftoverAfter = after
		}
		c.Deltas = append(c.Deltas, Delta{Field: field, Before: before, After: after, Decibel: decibel})
	}
	add("wavelength_m", a.LambdaM, b.LambdaM, false)
	add("fspl_db", a.FSPLDB, b.FSPLDB, true)
	add("eirp_dbm", a.EIRPDBm, b.EIRPDBm, true)
	add("received_power_dbm", a.PrDBm, b.PrDBm, true)
	if a.Assessment != nil && b.Assessment != nil {
		add("snr_db", a.Assessment.SNRDB, b.Assessment.SNRDB, true)
		add("margin_db", a.Assessment.MarginDB, b.Assessment.MarginDB, true)
	}
	return c
}
