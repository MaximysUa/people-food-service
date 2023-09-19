package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net"
	"net/http"
	"people-food-service/iternal/config"
	mwlogger "people-food-service/iternal/middleware/logger"
	person "people-food-service/iternal/person/db"
	ph "people-food-service/iternal/person/handlers"
	"people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"time"
)

func main() {
	cfg := config.GetConfig()

	logging.Init(cfg)
	logger := logging.GetLogger()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(logger))
	router.Use(middleware.Recoverer)
	//router.Use(middleware.URLFormat)

	client, err := postgresql.NewClient(context.TODO(), 5, cfg.Storage)
	defer client.Close()
	if err != nil {
		logger.Fatal(err)
	}
	repository := person.NewRepository(client, logger)
	router.Get("/api/person", ph.GetOne(logger, repository, context.TODO()))
	router.Get("/api/people", ph.GetList(logger, repository, context.TODO()))
	listener, listenErr := net.Listen("tcp", cfg.Listen.Port)
	logger.Infof("server is listening port %s", cfg.Listen.Port)

	if listenErr != nil {
		logger.Fatal(listenErr)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))

}
