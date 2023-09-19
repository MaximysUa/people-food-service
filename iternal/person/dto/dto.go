package persondto

import (
	"people-food-service/iternal/food"
	"people-food-service/iternal/person"
)

type ResponseDTO struct {
	Person         []person.Person `yaml:"person"`
	ResponseStatus string          `yaml:"response-status"`
}
type RequestDTO struct {
	UUID       string      `json:"uuid,omitempty"`
	Name       string      `json:"name" validate:"alphaunicode"`
	FamilyName string      `json:"family_name" validate:"alphaunicode"`
	Food       []food.Food `json:"food,omitempty"`
}
