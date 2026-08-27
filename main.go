package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"friis-link/internal/budget"
	"friis-link/internal/server"
)

const usage = `friis-link: free-space link budget.

A scenario JSON file gives frequency, distance, transmit power and the
antenna gains; the budget prints wavelength, FSPL, received power and,
when a noise section is present, SNR and the fading margin against the
required SNR.

usage:
  friis-link budget [-json] <scenario.json>
  friis-link reverse <scenario.json>
  friis-link compare <before.json> <after.json>
  friis-link help

scenario.json fields:
  frequency_hz    carrier frequency in hertz
  distance_km     slant range in kilometres
  tx_power_dbm    transmit power in dBm
  tx_gain_dbi     transmit antenna gain in dBi
  rx_gain_dbi     receive antenna gain in dBi
  extra_loss_db   optional additional loss in dB (polarisation, pointing)
  noise           optional object:
                    temperature_k    system temperature in kelvin
                    bandwidth_hz     receiver bandwidth in hertz
                    noise_figure_db  noise figure in dB
                    min_snr_db       required SNR in dB

Illegal inputs (non-positive distance or frequency, negative linear
gain, zero bandwidth, malformed JSON) are reported on stderr and exit
non-zero.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "budget":
		err = runBudget(os.Args[2:])
	case "-http", "--http":
		err = runHTTP(os.Args[2:])
	case "reverse":
		err = runReverse(os.Args[2:])
	case "compare":
		err = runCompare(os.Args[2:])
	case "help", "-h", "--help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "friis-link: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "friis-link: %v\n", err)
		os.Exit(1)
	}
}

func loadScenario(fs *flag.FlagSet) (*budget.Config, error) {
	if fs.NArg() != 1 {
		return nil, errors.New("command needs exactly one scenario JSON file")
	}
	return budget.LoadConfig(fs.Arg(0))
}

func runBudget(args []string) error {
	fs := flag.NewFlagSet("budget", flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "print the result as JSON instead of text")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := loadScenario(fs)
	if err != nil {
		return err
	}
	res, err := budget.Compute(cfg)
	if err != nil {
		return err
	}
	if *asJSON {
		data, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	fmt.Print(res.Report())
	return nil
}

func runReverse(args []string) error {
	fs := flag.NewFlagSet("reverse", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := loadScenario(fs)
	if err != nil {
		return err
	}
	res, err := budget.ReverseCompute(cfg)
	if err != nil {
		return err
	}
	fmt.Print(res.Report())
	return nil
}

func runCompare(args []string) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return errors.New("compare needs exactly two scenario JSON files")
	}
	before, err := budget.LoadConfig(fs.Arg(0))
	if err != nil {
		return err
	}
	after, err := budget.LoadConfig(fs.Arg(1))
	if err != nil {
		return err
	}
	resBefore, err := budget.Compute(before)
	if err != nil {
		return err
	}
	resAfter, err := budget.Compute(after)
	if err != nil {
		return err
	}
	fmt.Print(budget.Compare(resBefore, resAfter).Report())
	return nil
}

func runHTTP(args []string) error {
	addr := ":8080"
	if len(args) > 0 {
		addr = args[0]
	}
	if len(args) > 1 {
		return fmt.Errorf("-http accepts at most one listen address")
	}
	fmt.Printf("friis-link HTTP on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, server.New("web", "example/sband-10km.json"))
}
