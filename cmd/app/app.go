package main

import (
	"context"
	"fmt"
	"people-food-service/iternal/config"
	person "people-food-service/iternal/person/db"
	"people-food-service/pkg/client/postgresql"
)

func main() {
	//router := httprouter.New()
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 5, cfg.Storage)
	if err != nil {
		return
	}
	repository := person.NewRepository(client)
	all, err := repository.FindOne(context.TODO(), "Василий", "Уткин")
	fmt.Println(all)
}
