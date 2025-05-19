package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Teneieiza/Golang-API/models"
	"github.com/Teneieiza/Golang-API/services"
)

func HandlePetFoods(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pets, err := services.GetAllPetFoods()
		if err != nil {
			http.Error(w, "Cannot get pet foods", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(pets)

	case http.MethodPost:
		var pet models.PetFood
		if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		id, err := services.InsertPetFood(pet)
		if err != nil {
			http.Error(w, "Failed to insert pet food", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Pet food created with ID: %d", id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}




func HandlePetFood(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		pet, err := services.GetPetFoodById(id)
		if err != nil {
			http.Error(w, "Error retrieving pet food", http.StatusInternalServerError)
			return
		}
		if pet == nil {
			http.Error(w, "Pet food not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(pet)

	case http.MethodDelete:
		err := services.RemovePetFoodById(id)
		if err != nil {
			http.Error(w, "Error deleting pet food", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Pet food %d was deleted", id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
