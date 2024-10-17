package apiserver

import "net/http"

// responseWriter is a standart type, but with a status code.
type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader is a standart function for responseWriter, but it writes a status code into .code.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
