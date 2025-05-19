package models

type PetFood struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	AnimalType string  `json:"animaltype"`
	FoodType   string  `json:"foodtype"`
	ImageUrl   string  `json:"imageurl"`
}
