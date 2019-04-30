package session

import "net/http"

type headersBuffer struct {
	headers http.Header
}

func (w headersBuffer) Header() http.Header {
	return w.headers
}

func (w headersBuffer) Write(b []byte) (int, error) {
	return len(b), nil
}

func (w headersBuffer) WriteHeader(statusCode int) {}

func newHeadersBuffer() *headersBuffer {
	return &headersBuffer{
		headers: make(http.Header),
	}
}
