package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	cfg    *Config
	mux    *chi.Mux
	logger *slog.Logger
}

func NewServer(cfg *Config, logger *slog.Logger) (*HttpServer, error) {

	// create router
	mux := chi.NewRouter()

	// setup cors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowOriginFunc:  cfg.Cors.AllowOriginFunc,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		MaxAge:           cfg.Cors.MaxAge,
	}))

	// return server
	return &HttpServer{
		cfg:    cfg,
		mux:    mux,
		logger: logger,
	}, nil

}

func (s *HttpServer) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		s.mux.Use(middleware)
	}
}

func (s *HttpServer) Method(method string, pattern string, handler http.Handler) {
	s.mux.Method(method, pattern, handler)
}

func (s *HttpServer) Mount(pattern string, handler http.Handler) {
	s.mux.Mount(pattern, handler)
}

func (s *HttpServer) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.Address, s.mux)
}

func (s *HttpServer) ListenAndServeTLS() error {
	return http.ListenAndServeTLS(s.cfg.AddressTLS, s.cfg.CertFile, s.cfg.KeyFile, s.mux)
}
