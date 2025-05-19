package routes

import (
	"net/http"

	"github.com/Teneieiza/Golang-API/controllers"
)

func SetupRoutes() {
	http.HandleFunc("/api/petfood", controllers.HandlePetFoods)
	http.HandleFunc("/api/petfood/", controllers.HandlePetFood)
}