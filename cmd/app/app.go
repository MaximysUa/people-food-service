package main

import (
	"context"
	"net"
	"net/http"
	"people-food-service/iternal/config"
	food "people-food-service/iternal/food/db"
	person "people-food-service/iternal/person/db"
	"people-food-service/iternal/router"
	"people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"time"
)

func main() {
	cfg := config.GetConfig()
	ctx := context.TODO()

	logging.Init(cfg)
	logger := logging.GetLogger()
	//TODO how and where create new dbs if they are not exist
	client, err := postgresql.NewClient(ctx, 5, cfg.Storage)

	pRep := person.NewRepository(client, logger)
	fRep := food.NewRepository(client, logger)
	defer client.Close()
	if err != nil {
		logger.Fatal(err)
	}

	rout := router.New(ctx, logger, pRep, fRep)

	listener, listenErr := net.Listen("tcp", cfg.Listen.Port)
	logger.Infof("server is listening port %s", cfg.Listen.Port)

	if listenErr != nil {
		logger.Fatal(listenErr)
	}
	server := &http.Server{
		Handler:      rout,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))

}
