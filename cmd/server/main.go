package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"golden-signals-playground/internal/server"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	addr := getEnv("ADDR", ":8080")
	s := server.NewState()
	r := chi.NewRouter()

	r.Mount("/", server.Routes(s))
	r.Mount("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
