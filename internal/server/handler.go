package server

import (
	"net/http"
	"sync/atomic"
)

type State struct {
	LatencyMs int64
	FailPct   float64
	Ready     int32
	CpuMs     int64
}

func NewState() *State { return &State{LatencyMs: 50, FailPct: 0.00, Ready: 1, CpuMs: 0} }

func Routes(s *State) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/readyz", s.readyHandler)
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
