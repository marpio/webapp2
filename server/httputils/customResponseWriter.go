package httputils

import (
	"net/http"
)

type CustomResponseWriter interface {
	Header() http.Header
	Write(p []byte) (int, error)
	WriteHeader(status int)
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int

	// Written returns whether or not the ResponseWriter has been written.
	//Written() bool

	// Size returns the size of the response body.
	Size() int64
}

type customResponseWriter struct {
	responseWriter http.ResponseWriter
	status         int
	responseBytes  int64
}

func NewCustomResponseWriter(w http.ResponseWriter) CustomResponseWriter {
	return &customResponseWriter{
		w,
		200,
		0,
	}
}

func (rw *customResponseWriter) Header() http.Header {
	return rw.responseWriter.Header()
}

func (rw *customResponseWriter) Write(p []byte) (int, error) {
	written, err := rw.responseWriter.Write(p)
	rw.responseBytes += int64(written)
	return written, err
}

func (rw *customResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.responseWriter.WriteHeader(status)
}

func (rw *customResponseWriter) Status() int {
	return rw.status
}

func (rw *customResponseWriter) Size() int64 {
	return rw.responseBytes
}
