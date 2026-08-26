package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHealthAndBudget(t *testing.T) {
	root := testdataRoot(t)
	h := New(filepath.Join(root, "web"), filepath.Join(root, "example", "sband-10km.json"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("health %d", rr.Code)
	}
	body, err := os.ReadFile(filepath.Join(root, "example", "sband-10km.json"))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/budget", bytes.NewReader(body))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("budget %d %s", rr.Code, rr.Body.String())
	}
	var got BudgetResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.FSPLDB < 100 || got.FSPLDB > 140 {
		t.Fatalf("FSPL %g, want ~120 dB", got.FSPLDB)
	}
	if got.EIRPDBm != 33 {
		t.Fatalf("EIRP %g", got.EIRPDBm)
	}
}

func testdataRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "example", "sband-10km.json")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("cannot find example")
	return ""
}
