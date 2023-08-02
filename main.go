package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	apiCfg := apiConfig{}
	// Step 1:
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricInc(http.FileServer(http.Dir("."))))) // 3. Wrapping the FileServer with the MiddleWare we just wrote
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	// Step 2:
	corsMux := middlewareCors(mux)

	// Step 3:
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	// Step 4:
	// fmt.Println("Server listening on port 8080:")
	server.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler { // 2. write a new middleware method on a *apiConfig that increments the fileserverHits counter every time it's called
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++

		next.ServeHTTP(w, r) // Calling the next Handler in the Chain
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// writes the number of requests that have been counted as plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits)))
}
