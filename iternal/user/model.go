package user

import "people-food-service/iternal/food"

type User struct {
	UUID       string      `json:"uuid"`
	Name       string      `json:"name"`
	FamilyName string      `json:"family_name"`
	Food       []food.Food `json:"food"`
}
