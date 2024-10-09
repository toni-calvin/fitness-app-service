package seeder

import (
	"app/fitness-app-service/database"
	"app/fitness-app-service/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func SeedDatabaseExcercises() {
	database.DB.Exec("DROP TABLE IF EXISTS excercises CASCADE;")
	database.DB.Exec("ALTER SEQUENCE exercises_id_seq RESTART WITH 1")
	jsonData, err := os.ReadFile("seeder/data/exercises.json")
	if err != nil {
		log.Fatal("Failed to read JSON file: ", err)
	}

	var exercises []models.Exercise
	err = json.Unmarshal(jsonData, &exercises)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON: ", err)
	}

	for _, exercise := range exercises {

		if database.DB.Where("name = ?", exercise.Name).First(&models.Exercise{}).RowsAffected > 0 {
			break
		}

		newExercise := models.Exercise{
			Name:         exercise.Name,
			MuscleGroup:  exercise.MuscleGroup,
			MovementType: exercise.MovementType,
			Notes:        exercise.Notes,
		}
		database.DB.Create(&newExercise)
	}

	fmt.Println("Database seeded successfully!")
}
