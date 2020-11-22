package rest

import (
	"fmt"
	"net/http"
	"os"
)

// Decorator func
func decorator(f func()) func() {
	return func() {
		fmt.Println("before fn call")
		f()
		fmt.Println("after fn call")
	}
}

// BasicAuth decorator
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	// closure func block
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("MVP_USERNAME") && password == os.Getenv("MVP_PASSWORD") {
				next.ServeHTTP(w, r) // call request handler
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="invoice.mvp"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}
