package service

import "net/http"

type Signal interface {
	Stop()
	Start()
}

type MuxHandler struct {
	handler Handler

	// next stores the next.ServeHTTP to reduce memory allocate
	next func(rw http.ResponseWriter, r *http.Request)
}

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func (m MuxHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next)
}
