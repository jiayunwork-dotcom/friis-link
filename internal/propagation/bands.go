package propagation

import "fmt"

type Band struct {
	Name   string
	LowHz  float64
	HighHz float64
}

var StandardBands = []Band{
	{Name: "L", LowHz: 1e9, HighHz: 2e9},
	{Name: "S", LowHz: 2e9, HighHz: 4e9},
	{Name: "C", LowHz: 4e9, HighHz: 8e9},
	{Name: "X", LowHz: 8e9, HighHz: 12e9},
	{Name: "Ku", LowHz: 12e9, HighHz: 18e9},
	{Name: "K", LowHz: 18e9, HighHz: 27e9},
	{Name: "Ka", LowHz: 27e9, HighHz: 40e9},
}

func BandOf(frequencyHz float64) (Band, error) {
	if err := validateCommon(1000, frequencyHz); err != nil {
		return Band{}, err
	}
	for _, b := range StandardBands {
		if frequencyHz >= b.LowHz && frequencyHz < b.HighHz {
			return b, nil
		}
	}
	return Band{}, fmt.Errorf("frequency %g Hz is outside the standard band table", frequencyHz)
}

func BandName(frequencyHz float64) (string, error) {
	b, err := BandOf(frequencyHz)
	if err != nil {
		return "", err
	}
	return b.Name, nil
}

func IsBand(frequencyHz float64, name string) (bool, error) {
	b, err := BandOf(frequencyHz)
	if err != nil {
		return false, err
	}
	return b.Name == name, nil
}
