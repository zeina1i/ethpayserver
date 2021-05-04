package httputil

import (
	"net/http"
	"strconv"
)

type Server struct {
	Router *Router
	Config *Config
}

func NewServer(router *Router, config *Config) *Server {
	return &Server{Router: router, Config: config}
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":"+strconv.Itoa(s.Config.Port), s.Router)
}
