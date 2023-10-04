package fooddto

import (
	"people-food-service/iternal/food"
)

type RequestDTO struct {
	UUID  string  `json:"uuid,omitempty"`
	Name  string  `json:"name" validate:"alphaunicode"`
	Price float64 `json:"price"`
}

type ResponseDTO struct {
	Food           []food.Food `yaml:"food"`
	ResponseStatus string      `yaml:"response-status"`
	Err            string      `yaml:"err"`
}
