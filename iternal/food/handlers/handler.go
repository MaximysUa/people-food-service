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
	StatusOK  = "OK"
	StatusErr = "ERROR: "
)

// @Summary      GetOne
// @Description  get one food entity
// @Tags         Food
// @Security BasicAuth
// @Accept       json
// @Produce      json
// @Param        input body fooddto.RequestDTO true "food and price"
// @Success      200  {object}  fooddto.ResponseDTO
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/food [get]
func GetOne(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var res fooddto.ResponseDTO
		f, err := helper.ValidateFood(r)
		if err != nil {
			logger.Error(err)

			w.WriteHeader(http.StatusBadRequest)

			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}
		one, err := repos.FindOne(ctx, f.Name)
		if err != nil {
			logger.Errorf("failed to find food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}
		res.Food = append(res.Food, one)
		res.ResponseStatus = StatusOK

		render.JSON(w, r, res)
	}
}

// @Summary      GetList
// @Description  get list of food entity
// @Tags         Food
// @Security BasicAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  fooddto.ResponseDTO
// @Failure      400  {object}   error
// @Failure      500  {object}  error
// @Router       /api/food/all [get]
func GetList(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO

		all, err := repos.FindAll(ctx)
		if err != nil {
			logger.Errorf("failed to find all. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			res.ResponseStatus = StatusErr + fmt.Sprintf("failed to find all. Error: %v", err)
			render.JSON(w, r, res)

			return
		}
		res.Food = all
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}

// @Summary      Create
// @Description  create a food entity
// @Tags         Food
// @Security BasicAuth
// @Accept       json
// @Produce      json
// @Param        input body fooddto.RequestDTO true "food and price"
// @Success      201  {object}  fooddto.ResponseDTO
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/food [post]
func Create(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		req, err := helper.ValidateFood(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}
		if req.UUID != "" {
			err = validator.New().Var(req.UUID, "uuid")
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				res.ResponseStatus = StatusErr + err.Error()
				render.JSON(w, r, res)
				return
			}
		}
		uuid, err := repos.Create(ctx, food.Food(req))
		if err != nil && uuid == "" {
			logger.Errorf("failed to create food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			res.ResponseStatus = StatusErr + fmt.Sprintf("failed to create food. Error: %v", err)
			render.JSON(w, r, res)

			return
		} else if err != nil && uuid != "" {
			logger.Errorf("failed to create food. Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			res.ResponseStatus = StatusErr + fmt.Sprintf("food already exists. uuid: %s", uuid)
			render.JSON(w, r, res)

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

// @Summary      Delete
// @Description  delete a food entity
// @Tags         Food
// @Security BasicAuth
// @Accept       json
// @Produce      json
// @Param        input body fooddto.RequestDTO true "food and price"
// @Success      200  {object}  fooddto.ResponseDTO
// @Failure      400  {object}   error
// @Failure      500  {object}   error
// @Router       /api/food [delete]
func Delete(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		req, err := helper.ValidateFoodUUID(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}

		err = repos.Delete(ctx, food.Food(req))
		if err != nil {
			logger.Errorf("failed to delete food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			res.ResponseStatus = StatusErr + fmt.Sprintf("failed to delete food. Error: %v", err)
			render.JSON(w, r, res)

			return
		}
		w.WriteHeader(http.StatusOK)
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}

// @Summary      Update
// @Description  Update a food entity
// @Tags         Food
// @Security BasicAuth
// @Accept       json
// @Produce      json
// @Param        input body fooddto.RequestDTO true "food and price"
// @Success      200  {object}  fooddto.ResponseDTO
// @Failure      400  {object}   error
// @Failure      500  {object}   error
// @Router       /api/food [patch]
func Update(ctx context.Context, logger *logging.Logger, repos food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res fooddto.ResponseDTO
		req, err := helper.ValidateFood(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			res.ResponseStatus = StatusErr + err.Error()
			render.JSON(w, r, res)
			return
		}
		if req.UUID == "" {
			logger.Errorln("field ID is required")
			w.WriteHeader(http.StatusBadRequest)
			res.ResponseStatus = StatusErr + "field ID is required"
			render.JSON(w, r, res)

			return
		}
		err = repos.Update(ctx, food.Food(req))
		if err != nil {
			logger.Errorf("failed to update food. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			res.ResponseStatus = StatusErr + fmt.Sprintf("failed to update food. Error: %v", err)
			render.JSON(w, r, res)

			return
		}
		res.ResponseStatus = StatusOK
		render.JSON(w, r, res)
	}
}
