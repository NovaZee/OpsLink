package signal

import "net/http"

type MuxHandler struct {
	Handler Handler

	// Next stores the next.ServeHTTP to reduce memory allocate
	Next func(rw http.ResponseWriter, r *http.Request)
}

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func (m MuxHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.Handler.ServeHTTP(rw, r, m.Next)
}

type MessageSink interface {
}
