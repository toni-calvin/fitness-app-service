package database

import (
	"app/fitness-app-service/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	dsn := "host=localhost user=postgres password=notsecurepassword dbname=fitness_db port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	err = DB.AutoMigrate(&models.Mesocycle{}, &models.Microcycle{}, &models.Training{}, &models.Exercise{}, &models.Set{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}

func CleanTestDatabase(db *gorm.DB) {
	db.Exec("DROP TABLE IF EXISTS trainings CASCADE;")
	db.Exec("DROP TABLE IF EXISTS excercises CASCADE;")
	db.Exec("DROP TABLE IF EXISTS mesocycles CASCADE;")
	db.Exec("DROP TABLE IF EXISTS sets CASCADE;")
	db.Exec("DROP TABLE IF EXISTS training_excercises CASCADE;")
}
