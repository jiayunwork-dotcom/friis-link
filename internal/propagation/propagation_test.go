package propagation

import (
	"math"
	"testing"
)

const (
	testFreqHz = 2.45e9
	testDistM  = 10000.0
)

func TestRejectZeroDistance(t *testing.T) {
	if _, err := FSPLLinear(0, testFreqHz); err == nil {
		t.Error("FSPLLinear(0, f) = nil error, want rejection")
	}
	if _, err := FSPLdB(0, testFreqHz); err == nil {
		t.Error("FSPLdB(0, f) = nil error, want rejection")
	}
	if _, err := ReceivedPowerdBm(30, 3, 3, 0, testFreqHz, 0); err == nil {
		t.Error("ReceivedPowerdBm with zero distance = nil error, want rejection")
	}
}

func TestRejectNonPositiveFrequency(t *testing.T) {
	for _, f := range []float64{0, -1e9} {
		if _, err := Wavelength(f); err == nil {
			t.Errorf("Wavelength(%g) = nil error, want rejection", f)
		}
		if _, err := FSPLdB(testDistM, f); err == nil {
			t.Errorf("FSPLdB(d, %g) = nil error, want rejection", f)
		}
	}
}

func TestRejectNegativeGain(t *testing.T) {
	_, err := ReceivedPowerLinear(1, -2, 3, testDistM, testFreqHz)
	if err == nil {
		t.Error("negative tx gain accepted, want rejection")
	}
	_, err = ReceivedPowerLinear(1, 2, -3, testDistM, testFreqHz)
	if err == nil {
		t.Error("negative rx gain accepted, want rejection")
	}
}

func TestDistanceDoublingPrQuarter(t *testing.T) {
	pr1, err := ReceivedPowerLinear(1, 1, 1, testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("ReceivedPowerLinear(d) returned error: %v", err)
	}
	pr2, err := ReceivedPowerLinear(1, 1, 1, 2*testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("ReceivedPowerLinear(2d) returned error: %v", err)
	}
	ratio := pr1 / pr2
	if math.Abs(ratio-4) > 1e-9 {
		t.Errorf("Pr(d)/Pr(2d) = %g, want 4 (Pr scales as 1/d^2)", ratio)
	}

	db1, err := ReceivedPowerdBm(0, 0, 0, testDistM, testFreqHz, 0)
	if err != nil {
		t.Fatalf("ReceivedPowerdBm(d) returned error: %v", err)
	}
	db2, err := ReceivedPowerdBm(0, 0, 0, 2*testDistM, testFreqHz, 0)
	if err != nil {
		t.Fatalf("ReceivedPowerdBm(2d) returned error: %v", err)
	}
	delta := db2 - db1
	want := -20 * math.Log10(2) // exact 6.0206 dB, rounded "6 dB" in prose
	if math.Abs(delta-want) > 1e-9 {
		t.Errorf("Pr(2d)-Pr(d) = %.6f dB, want %.6f dB", delta, want)
	}
}

func TestFrequencyDoublingFSPLPlus6(t *testing.T) {
	lambda1, err := Wavelength(testFreqHz)
	if err != nil {
		t.Fatalf("Wavelength(f) returned error: %v", err)
	}
	lambda2, err := Wavelength(2 * testFreqHz)
	if err != nil {
		t.Fatalf("Wavelength(2f) returned error: %v", err)
	}
	if math.Abs(lambda2*2-lambda1) > 1e-9 {
		t.Errorf("lambda(2f) = %g, want half of lambda(f) = %g", lambda2, lambda1)
	}

	fspl1, err := FSPLdB(testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("FSPLdB(f) returned error: %v", err)
	}
	fspl2, err := FSPLdB(testDistM, 2*testFreqHz)
	if err != nil {
		t.Fatalf("FSPLdB(2f) returned error: %v", err)
	}
	delta := fspl2 - fspl1
	want := 20 * math.Log10(2) // exact 6.0206 dB
	if math.Abs(delta-want) > 1e-9 {
		t.Errorf("FSPL(2f)-FSPL(f) = %.6f dB, want %.6f dB", delta, want)
	}
}

func TestEIRPEqualsPtPlusGt(t *testing.T) {
	got := EIRPdBm(20, 15)
	if math.Abs(got-35) > 1e-12 {
		t.Errorf("EIRPdBm(20,15) = %g dBm, want 35 dBm", got)
	}
}

func TestGainPlus3DBPrPlus3DB(t *testing.T) {
	pr1, err := ReceivedPowerdBm(10, 3, 3, testDistM, testFreqHz, 0)
	if err != nil {
		t.Fatalf("baseline returned error: %v", err)
	}
	pr2, err := ReceivedPowerdBm(10, 6, 3, testDistM, testFreqHz, 0)
	if err != nil {
		t.Fatalf("boosted returned error: %v", err)
	}
	delta := pr2 - pr1
	if math.Abs(delta-3) > 1e-9 {
		t.Errorf("Pr(Gt+3dB)-Pr(Gt) = %.6f dB, want +3 dB", delta)
	}
}

func TestFSPLMatchesClosedForm(t *testing.T) {
	direct, err := FSPLdB(testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("FSPLdB returned error: %v", err)
	}
	closed, err := FSPLClosedForm(testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("FSPLClosedForm returned error: %v", err)
	}
	if math.Abs(direct-closed) > 1e-9 {
		t.Errorf("FSPLdB = %g, closed form = %g, want agreement", direct, closed)
	}
}

func TestWavelengthValue(t *testing.T) {
	lambda, err := Wavelength(testFreqHz)
	if err != nil {
		t.Fatalf("Wavelength returned error: %v", err)
	}
	want := 2.99792458e8 / testFreqHz
	if math.Abs(lambda-want) > 1e-12 {
		t.Errorf("lambda = %g, want %g", lambda, want)
	}
}

func TestExtraLossAppliedOnce(t *testing.T) {
	l, err := NewExtraLoss(2)
	if err != nil {
		t.Fatalf("NewExtraLoss(2) returned error: %v", err)
	}
	applied := l.Apply(-80)
	if math.Abs(applied+82) > 1e-12 {
		t.Errorf("Apply(-80) = %g, want -82 dBm", applied)
	}
	if _, err := NewExtraLoss(-1); err == nil {
		t.Error("NewExtraLoss(-1) = nil error, want rejection")
	}
}

func TestMaximumDistanceInverse(t *testing.T) {
	// Forward then inverse: the maximum distance for a received power
	// that already holds at d must be at least d.
	ptW, gt, gr := 1e-3, 2.0, 3.0
	prW, err := ReceivedPowerLinear(ptW, gt, gr, testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("ReceivedPowerLinear returned error: %v", err)
	}
	maxD, err := MaximumDistanceM(ptW, gt, gr, prW, testFreqHz)
	if err != nil {
		t.Fatalf("MaximumDistanceM returned error: %v", err)
	}
	if math.Abs(maxD-testDistM) > 1e-6 {
		t.Errorf("max distance = %g m, want %g m (round trip)", maxD, testDistM)
	}

	// Halving the required power extends the range by sqrt(2).
	maxD2, err := MaximumDistanceM(ptW, gt, gr, prW/2, testFreqHz)
	if err != nil {
		t.Fatalf("MaximumDistanceM(half power) returned error: %v", err)
	}
	want := maxD * math.Sqrt(2)
	if math.Abs(maxD2-want) > 1e-6 {
		t.Errorf("max distance at half power = %g m, want %g m", maxD2, want)
	}
}

func TestRequiredTransmitPowerRoundTrip(t *testing.T) {
	ptW, gt, gr := 1e-3, 2.0, 3.0
	prW, err := ReceivedPowerLinear(ptW, gt, gr, testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("ReceivedPowerLinear returned error: %v", err)
	}
	needed, err := RequiredTransmitPowerW(prW, gt, gr, testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("RequiredTransmitPowerW returned error: %v", err)
	}
	if math.Abs(needed-ptW) > ptW*1e-9 {
		t.Errorf("required Pt = %g W, want %g W", needed, ptW)
	}
}

func TestBandOf(t *testing.T) {
	band, err := BandOf(2.45e9)
	if err != nil {
		t.Fatalf("BandOf(2.45e9) returned error: %v", err)
	}
	if band.Name != "S" {
		t.Errorf("BandOf(2.45e9) = %q band, want S", band.Name)
	}
	ok, err := IsBand(10e9, "X")
	if err != nil {
		t.Fatalf("IsBand(10e9) returned error: %v", err)
	}
	if !ok {
		t.Error("10 GHz must be in the X band")
	}
	if _, err := BandOf(0); err == nil {
		t.Error("BandOf(0) = nil error, want rejection")
	}
}

func TestDecomposeFriisAgreesWithDirect(t *testing.T) {
	direct, err := ReceivedPowerLinear(1e-3, 2, 3, testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("ReceivedPowerLinear returned error: %v", err)
	}
	dec, err := DecomposeFriis(1e-3, 2, 3, testDistM, testFreqHz)
	if err != nil {
		t.Fatalf("DecomposeFriis returned error: %v", err)
	}
	if math.Abs(dec.ReceivedPower-direct) > direct*1e-12 {
		t.Errorf("decomposed Pr = %g W, direct = %g W, want agreement", dec.ReceivedPower, direct)
	}
	// Power density must drop as 1/d^2: doubling distance quarters it.
	s1, err := PowerDensityWPerM2(1, testDistM)
	if err != nil {
		t.Fatalf("PowerDensityWPerM2 returned error: %v", err)
	}
	s2, err := PowerDensityWPerM2(1, 2*testDistM)
	if err != nil {
		t.Fatalf("PowerDensityWPerM2(2d) returned error: %v", err)
	}
	if math.Abs(s1/s2-4) > 1e-9 {
		t.Errorf("S(d)/S(2d) = %g, want 4", s1/s2)
	}
}
