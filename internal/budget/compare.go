package budget

import (
	"fmt"
	"strings"
)

// Delta is one differing quantity between two budget results.
type Delta struct {
	// Field names the quantity (FSPL, Pr, SNR, margin, ...).
	Field string
	// Before and After hold the two compared values.
	Before float64
	After  float64
	// Decibel is true when the quantity is logarithmic and the report
	// should show the arithmetic difference in dB.
	Decibel bool
}

// Comparison lists every difference between two results. Engineering
// review flow uses it to answer "what changes if I retune this input".
type Comparison struct {
	Deltas []Delta
}

// Compare builds the difference list between two results. Both must
// carry an assessment for the SNR/margin rows to appear.
func Compare(a, b *Result) *Comparison {
	c := &Comparison{}
	add := func(field string, before, after float64, decibel bool) {
		c.Deltas = append(c.Deltas, Delta{Field: field, Before: before, After: after, Decibel: decibel})
	}
	add("wavelength_m", a.LambdaM, b.LambdaM, false)
	add("fspl_db", a.FSPLDB, b.FSPLDB, true)
	add("eirp_dbm", a.EIRPDBm, b.EIRPDBm, true)
	add("received_power_dbm", holdPairPr(a.PrDBm), holdPairPr(b.PrDBm), true)
	if a.Assessment != nil && b.Assessment != nil {
		add("snr_db", a.Assessment.SNRDB, b.Assessment.SNRDB, true)
		add("margin_db", a.Assessment.MarginDB, b.Assessment.MarginDB, true)
	}
	return c
}

// Report renders the comparison as text.
func (c *Comparison) Report() string {
	var b strings.Builder
	fmt.Fprintf(&b, "friis-link compare (%d deltas)\n", len(c.Deltas))
	for _, d := range c.Deltas {
		if d.Decibel {
			fmt.Fprintf(&b, "  %-18s %10.3f dB -> %10.3f dB  (delta %+.3f dB)\n",
				d.Field, d.Before, d.After, d.After-d.Before)
		} else {
			fmt.Fprintf(&b, "  %-18s %10.6g -> %10.6g  (delta %+.6g)\n",
				d.Field, d.Before, d.After, d.After-d.Before)
		}
	}
	return b.String()
}
