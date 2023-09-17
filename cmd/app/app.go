package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net"
	"net/http"
	"people-food-service/iternal/config"
	person "people-food-service/iternal/person/db"
	ph "people-food-service/iternal/person/handlers"
	"people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"time"
)

func main() {
	cfg := config.GetConfig()

	logging.Init(cfg.Env)
	logger := logging.GetLogger()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	client, err := postgresql.NewClient(context.TODO(), 5, cfg.Storage)
	defer client.Close()
	if err != nil {
		return
	}
	repository := person.NewRepository(client)
	router.Get("api/person", ph.GetOne(logger, repository, context.TODO()))

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
	//repository.Delete(context.TODO(), "b46f9850-97f7-4d60-9e8f-88ae58d72906")
	//repository.Update(context.TODO(), person2.Person{
	//
	//	UUID:       "6ebc8437-c7e2-4811-898e-7a002abd44d4",
	//	Name:       "Yamal",
	//	FamilyName: "Pilkin",
	//	Food:       nil,
	//})
}
