package budget

import (
	"fmt"
	"strings"

	"friis-link/internal/model"
	"friis-link/internal/noise"
)

func (r *Result) Report() string {
	var b strings.Builder
	fmt.Fprintf(&b, "friis-link budget\n")
	fmt.Fprintf(&b, "  frequency   : %s (%s band)\n", model.FormatFrequency(r.FrequencyHz), r.Band)
	fmt.Fprintf(&b, "  distance    : %s\n", model.FormatDistance(r.DistanceM))
	fmt.Fprintf(&b, "  wavelength  : %.6f m\n", r.LambdaM)
	fmt.Fprintf(&b, "  FSPL        : %.3f dB (linear %s)\n", r.FSPLDB, model.FormatRatio(r.FSPLLinear))
	fmt.Fprintf(&b, "  EIRP        : %s\n", model.FormatPowerDBm(r.EIRPDBm))
	fmt.Fprintf(&b, "  Pr          : %s (%s)\n", model.FormatPowerDBm(r.PrDBm), model.FormatPowerWatts(r.PrWatts))

	if r.Assessment == nil {
		b.WriteString("  noise       : n/a (no noise section in input)\n")
		b.WriteString("  SNR         : n/a\n")
		b.WriteString("  margin      : n/a\n")
		b.WriteString("  link        : n/a\n")
		return b.String()
	}

	a := r.Assessment
	fmt.Fprintf(&b, "  noise floor : %s (kTBF)\n", model.FormatPowerDBm(r.NoiseFloorDBm))
	fmt.Fprintf(&b, "  SNR         : %.3f dB\n", a.SNRDB)
	fmt.Fprintf(&b, "  min SNR     : %.3f dB\n", a.MinSNRDB)
	fmt.Fprintf(&b, "  margin      : %.3f dB\n", a.MarginDB)
	fmt.Fprintf(&b, "  link        : %s (%s)\n", marginVerdict(a), a.MarginState)
	return b.String()
}

func marginVerdict(a *noise.Assessment) string {
	if !a.Feasible {
		return "insufficient"
	}
	return "sufficient"
}
