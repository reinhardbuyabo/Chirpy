# Chirpy
## Building a fully-fledged Web Server from Scratch

## 1. Setting Up

Step 1: Creating a new http.ServeMux.

```go
mux := http.NewServeMux() // func NewServeMux() *ServeMux
```

Function signature for the above function(Returns a pointer to a ServeMux struct):
```go
func NewServeMux() *ServeMux
```
NewServeMux allocates and returns a new ServeMux.

ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.

Step 2: Wrapping the request multiplixer in a custom middleware function that adds Cross Origin Resource Sharing headers to the response. The function is defined as follows: 

```go
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) // func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
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
Handler returns the handler to use for the given request, consulting r.Method, r.Host, and r.URL.Path.
HandleFunc registers the handler function for the given pattern.

Middleware is a (loosely defined) term for any software or service that enables the parts of a system to communicate and manage data. It is the software that handles communication between components and input/output, so developers can focus on the specific purpose of their application

Step 3: Create a new http.Server and use the corsMux as the handler.
```go
server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}
```

A Server defines parameters for running an HTTP server.

Step 4: Use the server's ListenAndServe method to start the server
```go
server.ListenAndServe() // func (srv *Server) ListenAndServe() error
```

## 2. Fileservers
A fileserver is a kind of simple web server that serves static files from the host machine.

### STEPS

1. Add the HTML code above to a file called index.html in the same root directory as your server
2. Use the 
```go
http.NewServeMux
```
.Handle() method to **add a handler** for the root path (/).

3. Use a standard 
```go
http.FileServer // func FileServer(root FileSystem) Handler
```
as **the handler**
FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root

4. Use
```go
http.Dir
```
to **convert a filepath**, (in our case a dot: . which indicates the current directory) to a directory **for the http.FileServer**.
Re-build and run your server
Test your server by visiting http://localhost:8080 in your browser
Run the tests in the window on the right