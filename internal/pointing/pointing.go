package pointing

import (
	"fmt"
	"math"
)

func HalfPowerBeamwidth(gainDBi float64) (float64, error) {
	if !(gainDBi > 0) {
		return 0, fmt.Errorf("pointing: gain must be > 0 dBi")
	}
	lin := math.Pow(10, gainDBi/10)
	bw := 70 / math.Sqrt(lin)
	if !(bw > 0) || math.IsNaN(bw) || math.IsInf(bw, 0) {
		return 0, fmt.Errorf("pointing: beamwidth is not finite")
	}
	return bw, nil
}

func PointingLoss(offsetDeg, beamwidthDeg float64) (float64, error) {
	if offsetDeg < 0 {
		return 0, fmt.Errorf("pointing: offset must be >= 0")
	}
	if !(beamwidthDeg > 0) {
		return 0, fmt.Errorf("pointing: beamwidth must be > 0")
	}
	x := offsetDeg / beamwidthDeg
	loss := 12 * x * x
	if loss < 0 || math.IsNaN(loss) || math.IsInf(loss, 0) {
		return 0, fmt.Errorf("pointing: loss is not finite")
	}
	return loss, nil
}

func PolarizationMismatch(angleDeg float64) (float64, error) {
	if angleDeg < 0 || angleDeg > 90 {
		return 0, fmt.Errorf("pointing: polarization angle must be in [0, 90]")
	}
	if angleDeg >= 90-1e-9 {
		return 0, fmt.Errorf("pointing: orthogonal polarization is a total mismatch")
	}
	c := math.Cos(angleDeg * math.Pi / 180)
	if c <= 1e-12 {
		return 0, fmt.Errorf("pointing: orthogonal polarization is a total mismatch")
	}
	return -10 * math.Log10(c*c), nil
}

func CombineOnce(parts ...float64) (float64, error) {
	sum := 0.0
	for i, p := range parts {
		if p < 0 {
			return 0, fmt.Errorf("pointing: extra loss %d is negative", i)
		}
		if math.IsNaN(p) || math.IsInf(p, 0) {
			return 0, fmt.Errorf("pointing: extra loss %d is not finite", i)
		}
		sum += p
	}
	return sum, nil
}

func FromGeometry(gainDBi, offsetDeg, polDeg float64) (float64, error) {
	bw, err := HalfPowerBeamwidth(gainDBi)
	if err != nil {
		return 0, err
	}
	pt, err := PointingLoss(offsetDeg, bw)
	if err != nil {
		return 0, err
	}
	pol, err := PolarizationMismatch(polDeg)
	if err != nil {
		return 0, err
	}
	return CombineOnce(pt, pol)
}

func NarrowerBeamHurtsMore(lowGain, highGain, offsetDeg float64) error {
	a, err := FromGeometry(lowGain, offsetDeg, 0)
	if err != nil {
		return err
	}
	b, err := FromGeometry(highGain, offsetDeg, 0)
	if err != nil {
		return err
	}
	if b <= a {
		return fmt.Errorf("pointing: higher gain should hurt more at the same offset: %g vs %g", b, a)
	}
	return nil
}

func ZeroOffsetNoLoss(gainDBi float64) error {
	bw, err := HalfPowerBeamwidth(gainDBi)
	if err != nil {
		return err
	}
	loss, err := PointingLoss(0, bw)
	if err != nil {
		return err
	}
	if loss != 0 {
		return fmt.Errorf("pointing: boresight should be 0 dB, got %g", loss)
	}
	return nil
}

func OrthogonalRejected() error {
	_, err := PolarizationMismatch(90)
	if err == nil {
		return fmt.Errorf("pointing: 90° mismatch must be rejected")
	}
	return nil
}
