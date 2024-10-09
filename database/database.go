package database

import (
	"app/fitness-app-service/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(postgres.Open("fitnessapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Exercise{}, &models.Training{}, &models.TrainingExercise{}, &models.Set{}, &models.Mesocycle{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}
