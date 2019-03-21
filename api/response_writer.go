package api

import (
	"encoding/json"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func wrapHandlerWithResponseWriter(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rw := newResponseWriter(w)
		wrappedHandler.ServeHTTP(rw, req)
	})
}

func getStatusFromResponseWriter(w http.ResponseWriter) int {
	if rw, ok := w.(*responseWriter); ok {
		return rw.statusCode
	}
	return -1
}

//Write to the response and with the status code
func Write(w http.ResponseWriter, status int, text string) {
	WriteBytes(w, status, []byte(text))
}

//WriteBytes to the response and with the status code
func WriteBytes(w http.ResponseWriter, status int, text []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(text)
}

//WriteJSON to the response and with the status code
func WriteJSON(w http.ResponseWriter, status int, i interface{}) {
	bts, err := json.Marshal(i)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteBytes(w, status, bts)
}
