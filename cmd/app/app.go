package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"people-food-service/iternal/config"
	food "people-food-service/iternal/food/db"
	person "people-food-service/iternal/person/db"
	"people-food-service/iternal/router"
	"people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"syscall"
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
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Handler:      rout,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		logger.Fatal(server.Serve(listener))
	}()
	logger.Infof("server started")

	<-done

	logger.Infof("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server %v", err)

		return
	}
}
