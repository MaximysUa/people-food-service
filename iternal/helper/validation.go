package helper

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	persondto "people-food-service/iternal/person/dto"
)

// Validation name and familyName
// It should be not empty and consist only alhpabet characters
func Validation(r *http.Request) (persondto.RequestDTO, error) {
	req, err := renderToDTO(r)
	if err != nil {
		return persondto.RequestDTO{}, err
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return req, fmt.Errorf("invalid request. Error: %v", validateErr)

	}
	return req, nil
}

// ValidationUUID check uuid field
// It should be not empty and consist uuid
func ValidationUUID(r *http.Request) (persondto.RequestDTO, error) {
	req, err := renderToDTO(r)
	if err != nil {
		return persondto.RequestDTO{}, err
	}
	err = validator.New().Var(req.UUID, "uuid")
	if err != nil {
		return persondto.RequestDTO{}, fmt.Errorf("invalid request. Error: %v", err)
	}
	return req, nil
}

// renderToDTO rendering request to personRequestDTO
func renderToDTO(r *http.Request) (persondto.RequestDTO, error) {
	var req persondto.RequestDTO
	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		return req, errors.New("request body is empty")

	}
	if err != nil {
		return req, fmt.Errorf("failed to decode request body. Error: %v", err)
	}
	return req, nil
}
