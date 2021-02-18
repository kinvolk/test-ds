package main

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"io"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	hasher     func(w http.ResponseWriter, r *http.Request)
	httpServer http.Server
}

var _ http.Handler = (*Server)(nil)

func NewServer(ctx context.Context, hash string) *Server {
	f := failsum
	switch hash {
	case "sha256":
		f = sha256sum
	case "md5":
		f = md5sum
	}
	server := &Server{
		hasher: f,
	}
	server.httpServer = http.Server{
		Handler: server,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	return server
}

func sha256sum(w http.ResponseWriter, r *http.Request) {
	h := sha256.New()
	if _, err := io.Copy(h, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.Sum(nil))
}

func md5sum(w http.ResponseWriter, r *http.Request) {
	h := md5.New()
	if _, err := io.Copy(h, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.Sum(nil))
}

func failsum(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) Run(port int) error {
	s.httpServer.Addr = net.JoinHostPort("", strconv.Itoa(port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.hasher(w, r)
}
