package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Teneieiza/Golang-API/config"
	"github.com/Teneieiza/Golang-API/models"
)

func GetAllPetFoods() ([]models.PetFood, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := config.DB.QueryContext(ctx, "SELECT Id, Name, Price, AnimalType, FoodType, ImageUrl FROM petfood")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	petFoods := []models.PetFood{}
	for rows.Next() {
		var p models.PetFood
		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.AnimalType, &p.FoodType, &p.ImageUrl)
		if err != nil {
			return nil, err
		}
		petFoods = append(petFoods, p)
	}
	return petFoods, nil
}

func InsertPetFood(petFood models.PetFood) (int64, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := config.DB.ExecContext(context, "INSERT INTO petfood (Name, Price, AnimalType, FoodType, ImageUrl) VALUES (?, ?, ?, ?, ?)",
		petFood.Name,
		petFood.Price,
		petFood.AnimalType,
		petFood.FoodType,
		petFood.ImageUrl)

	if err != nil {
		log.Println("Insert error:", err)
		return 0, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println("LastInsertId error:", err)
		return 0, err
	}

	return insertId, nil
}

func GetPetFoodById(Id int) (*models.PetFood, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := config.DB.QueryRowContext(context, "SELECT Id, Name, Price, AnimalType, FoodType, ImageUrl FROM petfood WHERE Id = ?", Id)
	PetFood := &models.PetFood{}
	err := row.Scan(&PetFood.Id, &PetFood.Name, &PetFood.Price, &PetFood.AnimalType, &PetFood.FoodType, &PetFood.ImageUrl)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return PetFood, nil
}

func RemovePetFoodById(Id int) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := config.DB.ExecContext(context, "DELETE FROM petfood WHERE Id = ?", Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
