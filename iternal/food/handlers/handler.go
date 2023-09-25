package fh

import (
	"context"
	"net/http"
	fooddto "people-food-service/iternal/food/dto"
	"people-food-service/iternal/person"
	logging "people-food-service/pkg/client/logger"
)

func GetOne(ctx context.Context, logger *logging.Logger, repos person.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _ fooddto.ResponseDTO

	}
}
