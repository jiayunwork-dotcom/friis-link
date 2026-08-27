package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"friis-link/internal/budget"
	"friis-link/internal/rain"
	"friis-link/internal/twoway"
)

const maxBodyBytes = 1 << 20

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
}

type BudgetResponse struct {
	LambdaM  float64 `json:"lambda_m"`
	FSPLDB   float64 `json:"fspl_db"`
	EIRPDBm  float64 `json:"eirp_dbm"`
	PrDBm    float64 `json:"pr_dbm"`
	HasNoise bool    `json:"has_noise"`
	SNRDB    float64 `json:"snr_db,omitempty"`
	MarginDB float64 `json:"margin_db,omitempty"`
}

func New(staticDir, examplePath string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/example", makeExampleHandler(examplePath))
	mux.HandleFunc("/api/budget", handleBudget)
	mux.HandleFunc("/api/twoway", handleTwoWay)
	mux.HandleFunc("/api/rain", handleRain)
	mux.Handle("/", fileServer(staticDir))
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	writeJSON(w, http.StatusOK, HealthResponse{OK: true, Service: "friis-link"})
}

func handleBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	cfg, err := budget.ParseConfig(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := budget.Compute(cfg)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := BudgetResponse{
		LambdaM: res.LambdaM, FSPLDB: res.FSPLDB, EIRPDBm: res.EIRPDBm, PrDBm: res.PrDBm,
		HasNoise: res.HasAssessment(),
	}
	if res.Assessment != nil {
		out.SNRDB = res.Assessment.SNRDB
		out.MarginDB = res.Assessment.MarginDB
	}
	writeJSON(w, http.StatusOK, out)
}

type TwoWayRequest struct {
	FrequencyHz float64 `json:"frequency_hz"`
	DistanceM   float64 `json:"distance_m"`
	Ht          float64 `json:"ht"`
	Hr          float64 `json:"hr"`
}

func handleTwoWay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req TwoWayRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	dc, err := twoway.CrossoverHz(req.Ht, req.Hr, req.FrequencyHz)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	far, err := twoway.InTwoRay(req.DistanceM, req.Ht, req.Hr, req.FrequencyHz)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pe, err := twoway.PlaneEarthLossDB(req.DistanceM, req.Ht, req.Hr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"crossover_m": dc, "two_ray": far, "plane_earth_db": pe})
}

type RainRequest struct {
	K        float64 `json:"k"`
	Alpha    float64 `json:"alpha"`
	RateMMH  float64 `json:"rate_mmh"`
	LengthKm float64 `json:"length_km"`
}

func handleRain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req RainRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err := rain.Path(req.K, req.Alpha, req.RateMMH, req.LengthKm)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"path_db": p})
}

func makeExampleHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "GET required")
			return
		}
		f, err := os.Open(path)
		if err != nil {
			writeError(w, http.StatusNotFound, "example unavailable: "+err.Error())
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("request body is not valid JSON: %v", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func fileServer(dir string) http.Handler {
	inner := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		inner.ServeHTTP(w, r)
	})
}
