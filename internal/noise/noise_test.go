package noise

import (
	"math"
	"testing"
)

func TestNoiseFloorKTB(t *testing.T) {
	// N = k*T*B*F with F = 1 (0 dB noise figure).
	got, err := NoiseFloorWatts(290, 1e6, 1)
	if err != nil {
		t.Fatalf("NoiseFloorWatts returned error: %v", err)
	}
	want := 1.380649e-23 * 290 * 1e6
	if math.Abs(got-want) > want*1e-12 {
		t.Errorf("N = %g W, want %g W (k*T*B)", got, want)
	}

	// The noise figure must scale the floor linearly: F = 2 doubles it.
	withF, err := NoiseFloorWatts(290, 1e6, 2)
	if err != nil {
		t.Fatalf("NoiseFloorWatts(F=2) returned error: %v", err)
	}
	if math.Abs(withF-got*2) > got*2*1e-12 {
		t.Errorf("N(F=2) = %g W, want %g W", withF, got*2)
	}
}

func TestBandwidthDoublingSNRMinus3(t *testing.T) {
	const prDBm = -90.0
	n1, err := NoiseFloordBm(290, 1e6, 1)
	if err != nil {
		t.Fatalf("NoiseFloordBm(B) returned error: %v", err)
	}
	n2, err := NoiseFloordBm(290, 2e6, 1)
	if err != nil {
		t.Fatalf("NoiseFloordBm(2B) returned error: %v", err)
	}
	linearRatio := math.Pow(10, (n2-n1)/10)
	if math.Abs(linearRatio-2) > 1e-9 {
		t.Errorf("N(2B)/N(B) linear = %g, want 2", linearRatio)
	}

	snr1 := SNRdB(prDBm, n1)
	snr2 := SNRdB(prDBm, n2)
	delta := snr2 - snr1
	want := -10 * math.Log10(2) // exact 3.0103 dB
	if math.Abs(delta-want) > 1e-9 {
		t.Errorf("SNR(2B)-SNR(B) = %.6f dB, want %.6f dB", delta, want)
	}
}

func TestMarginBelowZeroLinkFails(t *testing.T) {
	// SNR just at the threshold -> margin 0 -> feasible.
	a := Assess(-100, -120, 20)
	if a.MarginDB != 0 {
		t.Errorf("margin = %g, want 0", a.MarginDB)
	}
	if !a.Feasible {
		t.Error("margin 0 must be feasible")
	}

	// SNR below threshold -> negative margin -> link insufficient.
	b := Assess(-100, -120, 25)
	if b.MarginDB != -5 {
		t.Errorf("margin = %g, want -5", b.MarginDB)
	}
	if b.Feasible {
		t.Error("negative margin must be infeasible")
	}
	if b.MarginState != MarginInsufficient {
		t.Errorf("state = %v, want insufficient", b.MarginState)
	}

	// Positive margin above the warn band -> healthy.
	c := Assess(-100, -120, 10)
	if c.MarginDB != 10 {
		t.Errorf("margin = %g, want 10", c.MarginDB)
	}
	if c.MarginState != MarginHealthy {
		t.Errorf("state = %v, want healthy", c.MarginState)
	}
}

func TestNoiseFloorRejectsZeroBandwidth(t *testing.T) {
	if _, err := NoiseFloorWatts(290, 0, 1); err == nil {
		t.Error("zero bandwidth accepted, want rejection")
	}
	if _, err := NoiseFloorWatts(0, 1e6, 1); err == nil {
		t.Error("zero temperature accepted, want rejection")
	}
	if _, err := NoiseFloorWatts(290, 1e6, 0); err == nil {
		t.Error("zero noise figure accepted, want rejection")
	}
}

func TestSNRLinearDivision(t *testing.T) {
	prW, err := NoiseFloorWatts(290, 1e6, 4)
	if err != nil {
		t.Fatalf("noise setup returned error: %v", err)
	}
	nW := prW / 2
	snr, err := SNRLinear(prW, nW)
	if err != nil {
		t.Fatalf("SNRLinear returned error: %v", err)
	}
	if math.Abs(snr-2) > 1e-12 {
		t.Errorf("SNR = %g, want 2", snr)
	}
}

func TestRequiredPowerForMinSNR(t *testing.T) {
	req := RequiredPowerForMinSNR(-110, 12)
	if math.Abs(req+98) > 1e-12 {
		t.Errorf("required Pr = %g dBm, want -98 dBm", req)
	}
}

func TestCascadeNoiseFigure(t *testing.T) {
	// A single stage returns its own noise figure.
	single, err := CascadeNoiseFigure([]Stage{{Name: "amp", GainLinear: 10, NoiseFigure: 2}})
	if err != nil {
		t.Fatalf("CascadeNoiseFigure(single) returned error: %v", err)
	}
	if math.Abs(single-2) > 1e-12 {
		t.Errorf("cascade(single) = %g, want 2", single)
	}

	// Two stages with unit gain: F = F1 + F2 - 1.
	two, err := CascadeNoiseFigure([]Stage{
		{Name: "a", GainLinear: 1, NoiseFigure: 2},
		{Name: "b", GainLinear: 1, NoiseFigure: 3},
	})
	if err != nil {
		t.Fatalf("CascadeNoiseFigure(two) returned error: %v", err)
	}
	if math.Abs(two-4) > 1e-12 {
		t.Errorf("cascade(two) = %g, want 4", two)
	}

	// A high-gain first stage suppresses the second stage's contribution.
	withGain, err := CascadeNoiseFigure([]Stage{
		{Name: "lna", GainLinear: 100, NoiseFigure: 1.5},
		{Name: "mixer", GainLinear: 1, NoiseFigure: 8},
	})
	if err != nil {
		t.Fatalf("CascadeNoiseFigure(lna+mixer) returned error: %v", err)
	}
	want := 1.5 + (8-1)/100.0
	if math.Abs(withGain-want) > 1e-9 {
		t.Errorf("cascade(lna+mixer) = %g, want %g", withGain, want)
	}
}

func TestSensitivityFor(t *testing.T) {
	s := SensitivityFor(-110, 12)
	if math.Abs(s.RequiredPowerDBm+98) > 1e-12 {
		t.Errorf("required power = %g dBm, want -98 dBm", s.RequiredPowerDBm)
	}
	excess := LinkBudgetExcess(-90, s)
	if math.Abs(excess-8) > 1e-12 {
		t.Errorf("excess = %g dB, want 8 dB", excess)
	}
}

func TestSystemTemperature(t *testing.T) {
	tsys, err := SystemTemperature(50, 240)
	if err != nil {
		t.Fatalf("SystemTemperature returned error: %v", err)
	}
	if tsys != 290 {
		t.Errorf("T_sys = %g K, want 290 K", tsys)
	}
	if _, err := SystemTemperature(-1, 290); err == nil {
		t.Error("negative antenna temperature accepted, want rejection")
	}
}

func TestTemperatureFromNoiseFloorRoundTrip(t *testing.T) {
	// N = k*T*B*F with T=290, B=1e6, F=1; recover T from N.
	nW, err := NoiseFloorWatts(290, 1e6, 1)
	if err != nil {
		t.Fatalf("NoiseFloorWatts returned error: %v", err)
	}
	T, err := TemperatureFromNoiseFloor(nW, 1e6, 1)
	if err != nil {
		t.Fatalf("TemperatureFromNoiseFloor returned error: %v", err)
	}
	if math.Abs(T-290) > 1e-6 {
		t.Errorf("recovered T = %g K, want 290 K", T)
	}
}
