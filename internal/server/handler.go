package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type State struct {
	LatencyMs int64
	FailPct   float64
	Ready     int32
	CpuMs     int64
}

func NewState() *State { return &State{LatencyMs: 50, FailPct: 0.00, Ready: 1} }

func Routes(s *State) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/healthz", InstrumentedHandler("healthz", http.HandlerFunc(healthHandler)))
	mux.Handle("/readyz", InstrumentedHandler("readyz", http.HandlerFunc(s.readyHandler)))
	mux.Handle("/work", InstrumentedHandler("work", http.HandlerFunc(s.workHandler)))
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *State) readyHandler(w http.ResponseWriter, _ *http.Request) {
	if atomic.LoadInt32(&s.Ready) == 1 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
		return
	}
	http.Error(w, "not ready", http.StatusServiceUnavailable)
}

func (s *State) workHandler(w http.ResponseWriter, r *http.Request) {
	lat := getInt(r, "latency_ms", int(atomic.LoadInt64(&s.LatencyMs)))
	fail := getFloat(r, "fail", s.FailPct)

	time.Sleep(time.Duration(lat) * time.Millisecond)

	if fail > 0 && rand.Float64() < fail {
		http.Error(w, "synthetic failure", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":         true,
		"latency_ms": lat,
		"fail":       fail,
		"ts":         time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func getInt(r *http.Request, k string, d int) int {
	if v := r.URL.Query().Get(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}
func getFloat(r *http.Request, k string, d float64) float64 {
	if v := r.URL.Query().Get(k); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n
		}
	}
	return d
}
