package ph

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"people-food-service/iternal/helper"
	"people-food-service/iternal/person"
	persondto "people-food-service/iternal/person/dto"
	logging "people-food-service/pkg/client/logger"
)

const (
	StatusOK  = "OK"
	StatusErr = "ERROR: "
)

func GetOne(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var res persondto.ResponseDTO
		req, err := helper.ValidatePerson(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}
		one, err := repos.FindOne(ctx, req.Name, req.FamilyName)
		if err != nil {
			logger.Errorf("failed to find a person. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}
		res.Person = append(res.Person, one)
		res.ResponseStatus = StatusOK
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
			render.JSON(w, r, fmt.Sprintf("failed to find all. Error: %v", err))
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
		req, err := helper.ValidatePerson(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}
		if req.UUID != "" {
			err = validator.New().Var(req.UUID, "uuid")
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, err)
				return
			}
		}

		uuid, err := repos.Create(ctx, person.Person(req))
		if err != nil && uuid == "" {
			logger.Errorf("failed to create person. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to create person. Error: %v", err))
			return
		} else if err != nil && uuid != "" {
			logger.Debugf("failed to create person. Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, fmt.Sprintf("person already exists. uuid: %s", uuid))
			return
		}
		w.WriteHeader(http.StatusCreated)
		res.Person = append(res.Person, person.Person{
			UUID:       uuid,
			Name:       req.Name,
			FamilyName: req.FamilyName,
			Food:       req.Food,
		})
		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}
func Delete(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res persondto.ResponseDTO
		req, err := helper.ValidatePersonUUID(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		err = repos.Delete(ctx, person.Person(req))
		if err != nil {
			logger.Errorf("failed to delete a person. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to delete a person. Error: %v", err))
			return
		}

		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}

func Update(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res persondto.ResponseDTO
		req, err := helper.ValidatePerson(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}
		if req.UUID == "" {
			logger.Errorln("Field ID is required")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "Field ID is required")
			return
		}

		err = repos.Update(ctx, person.Person(req))
		if err != nil {
			logger.Errorf("failed to update person. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to update person. Error: %v", err))
			return
		}
		res.Person = append(res.Person, person.Person(req))
		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}
