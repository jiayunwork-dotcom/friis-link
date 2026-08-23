package budget

import (
	"math"
	"path/filepath"
	"strings"
	"testing"

	"friis-link/internal/propagation"
)

var examplePath = filepath.Join("..", "..", "example", "sband-10km.json")

func loadExample(t *testing.T) *Config {
	t.Helper()
	cfg, err := LoadConfig(examplePath)
	if err != nil {
		t.Fatalf("LoadConfig(%s) returned error: %v", examplePath, err)
	}
	return cfg
}

func TestExampleSbandMagnitude(t *testing.T) {
	cfg := loadExample(t)
	res, err := Compute(cfg)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}

	if res.FSPLDB < 100 {
		t.Errorf("FSPL = %g dB, want at least 100 dB for an S-band 10 km link", res.FSPLDB)
	}

	closed, err := propagation.FSPLClosedForm(res.DistanceM, res.FrequencyHz)
	if err != nil {
		t.Fatalf("FSPLClosedForm returned error: %v", err)
	}
	if math.Abs(res.FSPLDB-closed) > 1e-9 {
		t.Errorf("FSPL = %g dB, closed form = %g dB, want agreement", res.FSPLDB, closed)
	}

	wantLambda := 2.99792458e8 / res.FrequencyHz
	if math.Abs(res.LambdaM-wantLambda) > 1e-12 {
		t.Errorf("lambda = %g m, want %g m", res.LambdaM, wantLambda)
	}
}

func TestExampleLinkVerdictSufficient(t *testing.T) {
	cfg := loadExample(t)
	res, err := Compute(cfg)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	if res.Assessment == nil {
		t.Fatal("assessment is nil, want a verdict for the example")
	}
	if !res.Assessment.Feasible {
		t.Errorf("example link must be feasible, margin = %g dB", res.Assessment.MarginDB)
	}
	if res.Assessment.MarginDB <= 0 {
		t.Errorf("example margin = %g dB, want positive", res.Assessment.MarginDB)
	}
}

func TestFrequencyTrendFSPL(t *testing.T) {
	// Doubling the frequency halves the wavelength, which raises FSPL by
	// exactly 6 dB. Because the noise floor is unchanged, SNR drops by
	// the same 6 dB and the feasibility verdict must follow: with the
	// required SNR set between the two SNR levels, the lower frequency
	// passes and the higher one fails. If FSPL were computed from
	// (4*pi*d*lambda) instead of (4*pi*d/lambda) the trend would invert
	// and the verdict would flip the wrong way.
	base := &Config{
		FrequencyHz: 2.45e9,
		DistanceKm:  10,
		TxPowerDBm:  0,
		TxGainDBi:   0,
		RxGainDBi:   0,
		Noise: &NoiseConfig{
			TemperatureK:  290,
			BandwidthHz:   1e6,
			NoiseFigureDB: 0,
			MinSNRDB:      -8,
		},
	}
	higher := *base
	higher.FrequencyHz = 2 * base.FrequencyHz

	low, err := Compute(base)
	if err != nil {
		t.Fatalf("Compute(f) returned error: %v", err)
	}
	high, err := Compute(&higher)
	if err != nil {
		t.Fatalf("Compute(2f) returned error: %v", err)
	}

	lambdaDelta := high.LambdaM / low.LambdaM
	if math.Abs(lambdaDelta-0.5) > 1e-9 {
		t.Errorf("lambda(2f)/lambda(f) = %g, want 0.5", lambdaDelta)
	}

	perOctave := 20 * math.Log10(2)
	fsplDelta := high.FSPLDB - low.FSPLDB
	if math.Abs(fsplDelta-perOctave) > 1e-9 {
		t.Errorf("FSPL(2f)-FSPL(f) = %.6f dB, want +%.6f dB", fsplDelta, perOctave)
	}

	snrDelta := high.Assessment.SNRDB - low.Assessment.SNRDB
	if math.Abs(snrDelta+perOctave) > 1e-9 {
		t.Errorf("SNR(2f)-SNR(f) = %.6f dB, want -%.6f dB", snrDelta, perOctave)
	}

	marginDelta := high.Assessment.MarginDB - low.Assessment.MarginDB
	if math.Abs(marginDelta+perOctave) > 1e-9 {
		t.Errorf("margin(2f)-margin(f) = %.6f dB, want -%.6f dB", marginDelta, perOctave)
	}

	if !low.Assessment.Feasible {
		t.Errorf("low-frequency link must pass, margin = %g dB", low.Assessment.MarginDB)
	}
	if high.Assessment.Feasible {
		t.Errorf("high-frequency link must fail, margin = %g dB", high.Assessment.MarginDB)
	}
}

func TestLinearAndDBPrAgree(t *testing.T) {
	cfg := loadExample(t)
	linearDBm, dbFormDBm, err := ComputeLinearAndDBResults(cfg)
	if err != nil {
		t.Fatalf("ComputeLinearAndDBResults returned error: %v", err)
	}
	if math.Abs(linearDBm-dbFormDBm) > 1e-6 {
		t.Errorf("linear Pr = %g dBm, decibel Pr = %g dBm, want agreement", linearDBm, dbFormDBm)
	}

	res, err := Compute(cfg)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	if math.Abs(res.CrossCheckDB) > 1e-6 {
		t.Errorf("cross check = %g dB, want near zero", res.CrossCheckDB)
	}
}

func TestBudgetRejectsInvalidInputs(t *testing.T) {
	validNoise := &NoiseConfig{TemperatureK: 290, BandwidthHz: 1e6, NoiseFigureDB: 2, MinSNRDB: 10}
	cases := []struct {
		name string
		cfg  Config
	}{
		{"zero distance", Config{FrequencyHz: 2.45e9, DistanceKm: 0, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3}},
		{"negative distance", Config{FrequencyHz: 2.45e9, DistanceKm: -1, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3}},
		{"zero frequency", Config{FrequencyHz: 0, DistanceKm: 10, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3}},
		{"negative frequency", Config{FrequencyHz: -1e9, DistanceKm: 10, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3}},
		{"negative extra loss", Config{FrequencyHz: 2.45e9, DistanceKm: 10, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3, ExtraLossDB: -0.5}},
		{"zero bandwidth", Config{FrequencyHz: 2.45e9, DistanceKm: 10, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3, Noise: &NoiseConfig{TemperatureK: 290, BandwidthHz: 0, NoiseFigureDB: 2, MinSNRDB: 10}}},
		{"zero temperature", Config{FrequencyHz: 2.45e9, DistanceKm: 10, TxPowerDBm: 30, TxGainDBi: 3, RxGainDBi: 3, Noise: &NoiseConfig{TemperatureK: 0, BandwidthHz: 1e6, NoiseFigureDB: 2, MinSNRDB: 10}}},
	}
	for _, c := range cases {
		if _, err := Compute(&c.cfg); err == nil {
			t.Errorf("%s: Compute accepted invalid config, want error", c.name)
		}
	}
	_ = validNoise
}

func TestBudgetParseInvalidJSON(t *testing.T) {
	inputs := []struct {
		name string
		data string
	}{
		{"malformed", `{`},
		{"wrong type", `{"frequency_hz": "2.45e9"}`},
		{"unknown field", `{"frequency_hz": 2.45e9, "bogus_field": 1}`},
	}
	for _, in := range inputs {
		if _, err := ParseConfig([]byte(in.data)); err == nil {
			t.Errorf("%s: ParseConfig accepted invalid JSON, want error", in.name)
		}
	}
}

func TestNoiseOptionalReport(t *testing.T) {
	cfg := &Config{
		FrequencyHz: 2.45e9,
		DistanceKm:  10,
		TxPowerDBm:  30,
		TxGainDBi:   3,
		RxGainDBi:   3,
	}
	res, err := Compute(cfg)
	if err != nil {
		t.Fatalf("Compute without noise returned error: %v", err)
	}
	if res.Assessment != nil {
		t.Error("assessment must be nil when noise is omitted")
	}
	report := res.Report()
	if !strings.Contains(report, "n/a") {
		t.Errorf("report without noise must mark SNR as n/a, got:\n%s", report)
	}
	if math.Abs(res.PrDBm+84.231) > 1e-3 {
		t.Errorf("Pr = %g dBm, want about -84.231 dBm", res.PrDBm)
	}
}

func TestReverseCompute(t *testing.T) {
	cfg := loadExample(t)
	rev, err := ReverseCompute(cfg)
	if err != nil {
		t.Fatalf("ReverseCompute returned error: %v", err)
	}
	// With Pt=30 dBm, Gt=Gr=3 dBi, the coverage at 2.45 GHz must be far
	// beyond the configured 10 km.
	if !rev.MaxDistanceFeasible {
		t.Errorf("10 km must be inside coverage, max distance = %g km", rev.MaxDistanceKm)
	}
	if rev.MaxDistanceKm <= 10 {
		t.Errorf("max distance = %g km, want more than 10 km", rev.MaxDistanceKm)
	}
	if rev.Sensitivity.RequiredPowerDBm >= -100 {
		t.Errorf("sensitivity = %g dBm, want below -100 dBm for this noise floor", rev.Sensitivity.RequiredPowerDBm)
	}
}

func TestJSONReport(t *testing.T) {
	cfg := loadExample(t)
	res, err := Compute(cfg)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	data, err := res.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON returned error: %v", err)
	}
	text := string(data)
	for _, key := range []string{"frequency_hz", "fspl_db", "snr_db", "margin_db", "link_feasible"} {
		if !strings.Contains(text, key) {
			t.Errorf("JSON output missing %q:\n%s", key, text)
		}
	}
}

func TestCompareDeltas(t *testing.T) {
	before := loadExample(t)
	after := *before
	after.TxPowerDBm = before.TxPowerDBm + 3
	resBefore, err := Compute(before)
	if err != nil {
		t.Fatalf("Compute(before) returned error: %v", err)
	}
	resAfter, err := Compute(&after)
	if err != nil {
		t.Fatalf("Compute(after) returned error: %v", err)
	}
	c := Compare(resBefore, resAfter)
	if len(c.Deltas) == 0 {
		t.Fatal("Compare produced no deltas")
	}
	// +3 dB transmit power must shift EIRP, Pr and SNR by +3 dB.
	for _, want := range []string{"eirp_dbm", "received_power_dbm", "snr_db"} {
		found := false
		for _, d := range c.Deltas {
			if d.Field == want {
				if d.After-d.Before != 3 {
					t.Errorf("%s delta = %g dB, want +3 dB", want, d.After-d.Before)
				}
				found = true
			}
		}
		if !found {
			t.Errorf("comparison missing %s", want)
		}
	}
}
