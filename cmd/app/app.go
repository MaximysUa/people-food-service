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

const (
	personURL = "/api/person"
	peopleURL = "/api/people"
)

func main() {
	cfg := config.GetConfig()
	ctx := context.TODO()

	logging.Init(cfg)
	logger := logging.GetLogger()

	client, err := postgresql.NewClient(ctx, 5, cfg.Storage)
	repository := person.NewRepository(client, logger)

	defer client.Close()
	if err != nil {
		logger.Fatal(err)
	}

	router := chi.NewRouter()
	//TODO мидлвари не работают why?
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(logger))
	router.Use(middleware.Recoverer)

	router.Post(personURL, ph.Create(ctx, logger, repository))
	router.Get(personURL, ph.GetOne(ctx, logger, repository))
	router.Get(peopleURL, ph.GetList(ctx, logger, repository))
	router.Delete(personURL, ph.Delete(ctx, logger, repository))
	router.Patch(personURL, ph.Update(ctx, logger, repository))

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
