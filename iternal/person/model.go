package person

import "people-food-service/iternal/food"

type Person struct {
	UUID       string      `json:"uuid"`
	Name       string      `json:"name"`
	FamilyName string      `json:"family_name"`
	Food       []food.Food `json:"food"`
}
