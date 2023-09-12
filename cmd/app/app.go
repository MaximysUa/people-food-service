package main

import (
	"context"
	"people-food-service/iternal/config"
	person2 "people-food-service/iternal/person"
	person "people-food-service/iternal/person/db"
	"people-food-service/pkg/client/postgresql"
)

func main() {
	//router := httprouter.New()
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 5, cfg.Storage)
	defer client.Close()
	if err != nil {
		return
	}
	repository := person.NewRepository(client)
	//repository.Delete(context.TODO(), "b46f9850-97f7-4d60-9e8f-88ae58d72906")
	repository.Update(context.TODO(), person2.Person{

		UUID:       "6ebc8437-c7e2-4811-898e-7a002abd44d4",
		Name:       "Yamal",
		FamilyName: "Pilkin",
		Food:       nil,
	})
}
