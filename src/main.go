package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	appEnv     = getEnv("APP_ENV", "local")
	appVersion = getEnv("APP_VERSION", "unknown")
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/version", handleVersion)
	return mux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from GitOps Demo! ENV=%s VERSION=%s\n", appEnv, appVersion)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"env":     appEnv,
		"version": appVersion,
	})
}

func main() {
	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on :%s (env=%s, version=%s)\n", port, appEnv, appVersion)
	if err := http.ListenAndServe(":"+port, newMux()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start server: %v\n", err)
		os.Exit(1)
	}
}
