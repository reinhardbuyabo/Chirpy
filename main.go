package main

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileServerHits int
	mutex          sync.Mutex
}

func main() {
	apiCfg := apiConfig{}
	r := chi.NewRouter()

	// Step 1:
	// mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	r.Handle("/app", apiCfg.middlewareMetricInc(fileServerHandler)) // You'll need to .Handle the fileserver handler twice, once for the /app/* path and once for the /app path (without the trailing slash).
	r.Handle("/app/*", apiCfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	// mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricInc(http.FileServer(http.Dir("."))))) // 3. Wrapping the FileServer with the MiddleWare we just wrote
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	// r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(http.StatusText(http.StatusOK)))
	// })

	// mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(http.StatusText(http.StatusOK)))
	// })
	adminRouter := chi.NewRouter()

	adminRouter.Get("/metrics", apiCfg.metricsHandler)

	// r.HandleFunc("/metrics", apiCfg.metricsHandler)

	// mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	// Step 2:

	corsMux := middlewareCors(r)
	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

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
	cfg.mutex.Lock()
	visits := cfg.fileServerHits
	cfg.mutex.Unlock()

	// Parse the HTML template from the file.
	tmpl, err := template.ParseFiles("admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response Content-Type header to "text/html".
	w.Header().Set("Content-Type", "text/html")

	// Execute the template with the PageData and write the output to the response.
	data := struct {
		Title  string
		Visits int
	}{
		Title:  "Welcome, Chirpy Admin",
		Visits: visits,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
