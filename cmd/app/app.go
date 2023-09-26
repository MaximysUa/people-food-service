package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net"
	"net/http"
	"people-food-service/iternal/config"
	food "people-food-service/iternal/food/db"
	fh "people-food-service/iternal/food/handlers"
	mwlogger "people-food-service/iternal/middleware/logger"
	person "people-food-service/iternal/person/db"
	ph "people-food-service/iternal/person/handlers"
	"people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"time"
)

const (
	personURL  = "/api/person"
	peopleURL  = "/api/people"
	foodURL    = "/api/food"
	allFoodURL = "/api/food/all"
)

func main() {
	cfg := config.GetConfig()
	ctx := context.TODO()

	logging.Init(cfg)
	logger := logging.GetLogger()

	client, err := postgresql.NewClient(ctx, 5, cfg.Storage)
	pRep := person.NewRepository(client, logger)
	fRep := food.NewRepository(client, logger)
	defer client.Close()
	if err != nil {
		logger.Fatal(err)
	}

	//TODO перенести создание роутера и регистрацию ручек в отдельный фаил
	router := chi.NewRouter()
	//TODO мидлвари не работают why?
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(logger))
	router.Use(middleware.Recoverer)

	router.Post(personURL, ph.Create(ctx, logger, pRep))
	router.Get(personURL, ph.GetOne(ctx, logger, pRep))
	router.Get(peopleURL, ph.GetList(ctx, logger, pRep))
	router.Delete(personURL, ph.Delete(ctx, logger, pRep))
	router.Patch(personURL, ph.Update(ctx, logger, pRep))

	router.Post(foodURL, fh.Create(ctx, logger, fRep))
	router.Get(foodURL, fh.GetOne(ctx, logger, fRep))
	router.Get(allFoodURL, fh.GetList(ctx, logger, fRep))
	router.Delete(foodURL, fh.Delete(ctx, logger, fRep))
	router.Patch(foodURL, fh.Update(ctx, logger, fRep))

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
