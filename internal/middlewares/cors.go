package middlewares

import (
	"log"
	"net/http"
)

const (
	options          string = "OPTIONS"
	allowOrigin      string = "Access-Control-Allow-Origin"
	allowMethods     string = "Access-Control-Allow-Methods"
	allowHeaders     string = "Access-Control-Allow-Headers"
	allowCredentials string = "Access-Control-Allow-Credentials"
	exposeHeaders    string = "Access-Control-Expose-Headers"
	credentials      string = "true"
	origin           string = "Origin"
	methods          string = "POST, GET, OPTIONS, PUT, DELETE, HEAD, PATCH"

	// If you want to expose some other headers add it here
	headers string = "Access-Control-Allow-Origin, Accept, Accept-Encoding, Authorization, Content-Length, Content-Type, X-CSRF-Token"
)

// AllowCORS handler will allow cross-origin HTTP requests
func AllowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set allow origin to match origin of our request or fall back to *
		if o := r.Header.Get(origin); o != "" {
			w.Header().Set(allowOrigin, o)
		} else {
			w.Header().Set(allowOrigin, "*")
		}

		// Set other headers
		w.Header().Set(allowHeaders, headers)
		w.Header().Set(allowMethods, methods)
		w.Header().Set(allowCredentials, credentials)
		w.Header().Set(exposeHeaders, headers)

		// If this was preflight options request let's write empty ok response and return
		if r.Method == options {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(nil); err != nil {
				log.Println("Error  write header")
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
