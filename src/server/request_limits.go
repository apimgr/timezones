package server

import "net/http"

const (
	maxRequestTargetBytes = 8192
	maxRequestBodyBytes   = 10 << 20
)

func requestLimitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.RequestURI()) > maxRequestTargetBytes {
			http.Error(w, "Request target is too long", http.StatusRequestURITooLong)
			return
		}
		if r.ContentLength > maxRequestBodyBytes {
			http.Error(w, "Request body is too large", http.StatusRequestEntityTooLarge)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		next.ServeHTTP(w, r)
	})
}
