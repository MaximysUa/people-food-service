package router

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"people-food-service/iternal/food"
	fh "people-food-service/iternal/food/handlers"
	mwlogger "people-food-service/iternal/middleware/logger"
	"people-food-service/iternal/person"
	ph "people-food-service/iternal/person/handlers"
	logging "people-food-service/pkg/client/logger"
)

const (
	personURL  = "/api/person"
	peopleURL  = "/api/people"
	foodURL    = "/api/food"
	allFoodURL = "/api/food/all"
)

func New(ctx context.Context, logger *logging.Logger,
	pRep person.Repository, fRep food.Repository) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	router.Use(mwlogger.New(logger))
	router.Use(middleware.Recoverer)

	router.Post(personURL, ph.Create(ctx, logger, pRep))
	router.Get(personURL, ph.GetOne(ctx, logger, pRep))
	router.Get(peopleURL, ph.GetList(ctx, logger, pRep))
	router.Delete(personURL, ph.Delete(ctx, logger, pRep))
	router.Patch(personURL, ph.Update(ctx, logger, pRep))

	router.Post(foodURL, fh.Create(ctx, logger, fRep))
	router.Get(foodURL, fh.GetOne(ctx, logger, fRep))
	router.Get(allFoodURL, fh.GetList(ctx, logger, fRep))
	router.Delete(foodURL, fh.Delete(ctx, logger, fRep))
	router.Patch(foodURL, fh.Update(ctx, logger, fRep))

	return router
}
