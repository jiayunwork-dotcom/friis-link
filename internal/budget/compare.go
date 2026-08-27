package budget

import (
	"fmt"
	"strings"
)

type Delta struct {
	Field   string
	Before  float64
	After   float64
	Decibel bool
}

type Comparison struct {
	Deltas []Delta
}

func Compare(a, b *Result) *Comparison {
	c := &Comparison{}
	add := func(field string, before, after float64, decibel bool) {
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
