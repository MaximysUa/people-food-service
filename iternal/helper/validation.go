package helper

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	fooddto "people-food-service/iternal/food/dto"
	persondto "people-food-service/iternal/person/dto"
)

// ValidatePerson name and familyName
// It should be not empty and consist only alhpabet characters
func ValidatePerson(r *http.Request) (persondto.RequestDTO, error) {
	req, err := renderToPersonDTO(r)
	if err != nil {
		return persondto.RequestDTO{}, err
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return req, fmt.Errorf("invalid request. Error: %v", validateErr)

	}
	return req, nil
}

// ValidatePersonUUID check uuid field
// It should be not empty and consist uuid
func ValidatePersonUUID(r *http.Request) (persondto.RequestDTO, error) {
	req, err := renderToPersonDTO(r)
	if err != nil {
		return persondto.RequestDTO{}, err
	}
	err = validator.New().Var(req.UUID, "uuid")
	if err != nil {
		return persondto.RequestDTO{}, fmt.Errorf("invalid request. Error: %v", err)
	}
	return req, nil
}

// renderToPersonDTO rendering request to personRequestDTO
func renderToPersonDTO(r *http.Request) (persondto.RequestDTO, error) {
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

// renderToFoodDTO rendering request to foodRequestDTO
func renderToFoodDTO(r *http.Request) (fooddto.RequestDTO, error) {
	var req fooddto.RequestDTO
	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		return req, errors.New("request body is empty")

	}
	if err != nil {
		return req, fmt.Errorf("failed to decode request body. Error: %v", err)
	}
	return req, nil
}

// ValidateFoodUUID check uuid field
// It should be not empty and consist uuid
func ValidateFoodUUID(r *http.Request) (fooddto.RequestDTO, error) {
	req, err := renderToFoodDTO(r)
	if err != nil {
		return fooddto.RequestDTO{}, err
	}
	err = validator.New().Var(req.UUID, "uuid")
	if err != nil {
		return fooddto.RequestDTO{}, fmt.Errorf("invalid request. Error: %v", err)
	}
	return req, nil
}

// ValidateFood name and price
// Name should be not empty and consists only alhpabet characters
// Price should be not empty and consists digits
func ValidateFood(r *http.Request) (fooddto.RequestDTO, error) {
	req, err := renderToFoodDTO(r)
	if err != nil {
		return fooddto.RequestDTO{}, err
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return fooddto.RequestDTO{}, fmt.Errorf("invalid request. Error: %v", validateErr)

	}
	if req.Price <= 0.0 {
		return fooddto.RequestDTO{}, fmt.Errorf("price have to be greater than 0")
	}

	return req, nil
}
