package ph

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"people-food-service/iternal/helper"
	"people-food-service/iternal/person"
	persondto "people-food-service/iternal/person/dto"
	logging "people-food-service/pkg/client/logger"
)

func GetOne(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var res persondto.ResponseDTO
		req, err := helper.Validation(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}
		one, err := repos.FindOne(ctx, req.Name, req.FamilyName)
		if err != nil {
			return
		}
		res.Person = append(res.Person, one)
		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}
func GetList(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res persondto.ResponseDTO

		all, err := repos.FindAll(ctx)
		if err != nil {
			logger.Errorf("failed to find all. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Errorf("failed to find all. Error: %v", err))
			return
		}
		res.Person = all
		res.ResponseStatus = "ok"
		render.JSON(w, r, res)
	}
}

func Create(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res persondto.ResponseDTO
		req, err := helper.Validation(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		err = repos.Create(ctx, person.Person(req))
		if err != nil {
			logger.Errorf("failed to create person. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Errorf("failed to create person. Error: %v", err))
			return
		}
		w.WriteHeader(http.StatusCreated)
		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}
func Delete(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res persondto.ResponseDTO
		req, err := helper.ValidationUUID(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		err = repos.Delete(ctx, person.Person(req))
		if err != nil {
			logger.Errorf("failed to delete person. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Errorf("failed to delete person. Error: %v", err))
			return
		}

		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}
