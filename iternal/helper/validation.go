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
// Validationt name and familyName
// It should be not empty and consist only alhpabet characters
func Validation(r *http.Request) (persondto.RequestDTO, error) {
	req, err := renderToDTO(r)
	if err != nil{
		return req, err	
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return req, fmt.Errorf("invalid request. Error: %v", validateErr)

	}
	return req, nil
}
//TODO write the validatoin of uuid
func ValidationUUID(r *http.Request) (persondto.RequestDTO, error){
	req, err := renderToDTO(r)
	if err != nil{
		return req, err	
	}

	return req, nil
}
func renderToDTO(r *http.Request) (persondto.RequestDTO, error){
	var req persondto.RequestDTO
	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		return req, errors.New("request body is empty")

	}
	if err != nil {

		return req, fmt.Errorf("failed to decode request body. Error: %v", err)
	}
}
