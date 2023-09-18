package apiserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"people-food-service/iternal/config"
	mwlogger "people-food-service/iternal/middleware/logger"
	logging "people-food-service/pkg/client/logger"
)

type APIServer struct {
	cfg    *config.Config
	logger *logging.Logger
	router *chi.Mux
}

func New(cfg *config.Config) *APIServer {
	logging.Init(cfg)
	return &APIServer{
		cfg:    cfg,
		logger: logging.GetLogger(),
		router: chi.NewRouter(),
	}
}
func (s *APIServer) Start() error {

	s.configureRouter()

	s.logger.Infof("Starting api server, port: %s", s.cfg.Listen.Port)
	return http.ListenAndServe(s.cfg.Listen.Port, s.router)
}
func (s *APIServer) configureRouter() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(mwlogger.New(s.logger))
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)
}

func (s *APIServer) Logger() *logging.Logger {
	return s.logger
}
func (s *APIServer) Router() *chi.Mux {
	return s.router
}
