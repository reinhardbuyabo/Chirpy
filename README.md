# Chirpy
## Building a fully-fledged Web Server from Scratch

## 1. Setting Up

Step 1: Creating a new http.ServeMux.

```go
mux := http.NewServeMux()
```

Function signature for the above function(Returns a pointer to a ServeMux struct):
```go
func NewServeMux() *ServeMux
```

ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.

Step 2: Wrapping the request multiplixer in a custom middleware function that adds Cross Origin Resource Sharing headers to the response. The function is defined as follows: 

```go
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
```

Middleware is a (loosely defined) term for any software or service that enables the parts of a system to communicate and manage data. It is the software that handles communication between components and input/output, so developers can focus on the specific purpose of their application

Step 3: Create a new http.Server and use the corsMux as the handler

Step 4: Use the server's ListenAndServe method to start the server



