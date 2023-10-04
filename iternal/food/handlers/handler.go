package fh

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"people-food-service/iternal/food"
	fooddto "people-food-service/iternal/food/dto"
	"people-food-service/iternal/helper"
	logging "people-food-service/pkg/client/logger"
)
const (
	StatusOK = "OK"
	)
func GetOne(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		f, err := helper.ValidateFood(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}
		one, err := repos.FindOne(ctx, f.Name)
		if err != nil {
			logger.Errorf("failed to find food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to find food. Error: %v", err))
			return
		}
		res.Food = append(res.Food, one)
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}
func GetList(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO

		all, err := repos.FindAll(ctx)
		if err != nil {
			logger.Errorf("failed to find all. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to find all. Error: %v", err))
			return
		}
		res.Food = all
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}
func Create(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		req, err := helper.ValidateFood(r)
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
		uuid, err := repos.Create(ctx, food.Food(req))
		if err != nil && uuid == "" {
			logger.Errorf("failed to create food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to create food. Error: %v", err))
			return
		} else if err != nil && uuid != "" {
			logger.Errorf("failed to create food. Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, fmt.Sprintf("food already exists. uuid: %s", uuid))
			return
		}
		w.WriteHeader(http.StatusCreated)
		res.Food = append(res.Food, food.Food{
			UUID:  uuid,
			Name:  req.Name,
			Price: req.Price,
		})
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}

func Delete(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		req, err := helper.ValidateFoodUUID(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		err = repos.Delete(ctx, food.Food(req))
		if err != nil {
			logger.Errorf("failed to delete food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to delete food. Error: %v", err))
			return
		}
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}
func Update(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		req, err := helper.ValidateFood(r)
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
		err = repos.Update(ctx, food.Food(req))
		if err != nil {
			logger.Errorf("failed to update food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("failed to update food. Error: %v", err))
			return
		}
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}
