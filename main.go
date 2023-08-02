package main

import (
	"net/http"
)

func main() {
	// Step 1:
	mux := http.NewServeMux().Handle("/", http.FileServer(http.Dir(".")))
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
