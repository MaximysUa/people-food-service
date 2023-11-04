package router

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	//_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	_ "people-food-service/cmd/app/docs"
	"people-food-service/iternal/config"
	"people-food-service/iternal/food"
	fh "people-food-service/iternal/food/handlers"
	mwlogger "people-food-service/iternal/middleware/logger"
	"people-food-service/iternal/person"
	ph "people-food-service/iternal/person/handlers"
	logging "people-food-service/pkg/client/logger"
)

const (
	personURL  = "/person"
	peopleURL  = "/people"
	foodURL    = "/food"
	allFoodURL = "/food/all"
)

func New(ctx context.Context, logger *logging.Logger,
	pRep person.Repository, fRep food.Repository, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	router.Use(mwlogger.New(logger))
	router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.BasicAuth("people-food-service", map[string]string{
			cfg.User: cfg.Password,
		}))
		r.Post(personURL, ph.Create(ctx, logger, pRep))
		r.Get(personURL, ph.GetOne(ctx, logger, pRep))
		r.Get(peopleURL, ph.GetList(ctx, logger, pRep))
		r.Delete(personURL, ph.Delete(ctx, logger, pRep))
		r.Patch(personURL, ph.Update(ctx, logger, pRep))

		r.Post(foodURL, fh.Create(ctx, logger, fRep))
		r.Get(foodURL, fh.GetOne(ctx, logger, fRep))
		r.Get(allFoodURL, fh.GetList(ctx, logger, fRep))
		r.Delete(foodURL, fh.Delete(ctx, logger, fRep))
		r.Patch(foodURL, fh.Update(ctx, logger, fRep))
	})
	logger.Info(router.Routes())
	return router
}
