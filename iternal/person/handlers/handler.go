package ph

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"people-food-service/iternal/food"
	"people-food-service/iternal/person"
	logging "people-food-service/pkg/client/logger"
)

type Response struct {
	Person         []person.Person `yaml:"person"`
	ResponseStatus string          `yaml:"response-status"`
}

type Request struct {
	UUID       string      `json:"uuid,omitempty"`
	Name       string      `json:"name"`
	FamilyName string      `json:"family_name"`
	Food       []food.Food `json:"food,omitempty"`
}

func GetOne(logger *logging.Logger, repos person.Repository, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var res Response
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")
			render.JSON(w, r, fmt.Errorf("request body is empty"))
			return
		}
		if err != nil {
			logger.Errorf("failed to decode request body. Error: %v", err)
			render.JSON(w, r, fmt.Errorf("failed to decode request body. Error: %v", err))
			return
		}
		logger.Tracef("request body decoded. Request: %v", req)
		one, err := repos.FindOne(ctx, req.Name, req.FamilyName)
		if err != nil {
			return
		}
		res.Person = append(res.Person, one)
		res.ResponseStatus = "Ok"
		render.JSON(w, r, res)
	}
}
