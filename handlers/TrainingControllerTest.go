package handlers

import (
	"app/fitness-app-service/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CleanTestDatabase(db *gorm.DB) {
	db.Exec("DROP TABLE IF EXISTS trainings CASCADE;")
	db.Exec("DROP TABLE IF EXISTS excercises CASCADE;")
	db.Exec("DROP TABLE IF EXISTS mesocycles CASCADE;")
	db.Exec("DROP TABLE IF EXISTS sets CASCADE;")
	db.Exec("DROP TABLE IF EXISTS training_excercises CASCADE;")
}

func SetupTestDatabase() *gorm.DB {
	fmt.Println("Setting up db")
	dsn := "host=localhost user=postgres password=notsecurepassword dbname=test_fitness_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		CleanTestDatabase(tx)
		if err := tx.AutoMigrate(&models.Mesocycle{}, &models.TrainingExercise{}, &models.Training{}, &models.Exercise{}, &models.Set{}); err != nil {
			return err
		}
		SeedTestDatabase(tx)
		return nil
	})

	if err != nil {
		fmt.Println("Error setting up database:", err)
		return nil
	}

	return db
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Pass db to handlers, e.g., with dependency injection
	r.GET("/training", func(c *gin.Context) {
		GetTrainings(c, db)
	})
	r.GET("/training/:id", func(c *gin.Context) {
		GetTrainingByID(c, db)
	})
	r.POST("/training", func(c *gin.Context) {
		CreateTraining(c, db)
	})
	r.PUT("/training/:id", func(c *gin.Context) {
		UpdateTraining(c, db)
	})
	r.DELETE("/training/:id", func(c *gin.Context) {
		DeleteTraining(c, db)
	})

	return r
}

func SeedTestDatabase(db *gorm.DB) *gorm.DB {
	mesocycle := models.Mesocycle{
		PreparationLevel: 5,
		Comments:         "Test comments",
		Objectives:       "Test objectives",
		Trainings:        []models.Training{},
	}

	Training := models.Training{
		MesocycleID: mesocycle.ID,
		Date:        "2024-10-03",
		TotalReps:   50,
		TotalWeight: 500.0,
		Excercises:  []models.TrainingExercise{},
	}

	mesocycle.Trainings = append(mesocycle.Trainings, Training)

	if err := db.Create(&mesocycle).Error; err != nil {
		fmt.Println("Error creating Mesocycle:", err)
		return nil
	}

	return db
}

func TestGetTrainings(t *testing.T) {
	db := SetupTestDatabase()
	router := SetupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/training", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Training
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Greater(t, len(response), 0) // Ensure there are training days in the response
}

func TestGetTraining(t *testing.T) {
	db := SetupTestDatabase()
	router := SetupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/training/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Training
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "2024-10-03", response.Date)
}

func TestAddTraining(t *testing.T) {
	db := SetupTestDatabase()
	router := SetupRouter(db)
	newDay := models.Training{
		MesocycleID: 1,
		Date:        "2024-10-10",
	}

	jsonBody, err := json.Marshal(newDay)
	if err != nil {
		fmt.Println("Error marshalling the JSON", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/training", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Training
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "2024-10-10", response.Date)
}

func TestUpdateTraining(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router
	updatedDay := models.Training{
		TotalReps: 0,
	}

	jsonBody, err := json.Marshal(updatedDay)
	if err != nil {
		fmt.Println("Error marshalling the JSON", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/training/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Training
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Zero(t, response.TotalReps)
	assert.Equal(t, "2024-10-03", response.Date)
}

func TestDeleteTraining(t *testing.T) {
	db := SetupTestDatabase()
	router := SetupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/training/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Training day deleted successfully", response["message"]) // Ensure the response confirms deletion
}
