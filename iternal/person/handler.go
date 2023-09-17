package person

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"people-food-service/iternal/handlers"
	logging "people-food-service/pkg/client/logger"
)

// TODO delete before prod
var _ handlers.Handler = &handler{}

const (
	personURL = "api/person"
	peopleURL = "api/people"
)

type handler struct {
	logger logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: *logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(personURL, h.GetPerson)
	router.GET(peopleURL, h.GetPeople)
	router.POST(personURL, h.CreatePerson)
	router.PUT(personURL, h.UpdatePerson)
	router.DELETE(personURL, h.DeletePerson)
}

func (h *handler) GetPeople(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (h *handler) GetPerson(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (h *handler) CreatePerson(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (h *handler) UpdatePerson(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (h *handler) DeletePerson(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
