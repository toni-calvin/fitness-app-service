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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDatabase() *gorm.DB {
	fmt.Println("Setting up db")
	db, _ := gorm.Open(sqlite.Open("../testfitnessapp.db"), &gorm.Config{})
	db.AutoMigrate(&models.Mesocycle{}, &models.TrainingDayExercise{}, &models.TrainingDay{}, &models.Exercise{}, &models.Set{})
	db = SeedTestDatabase(db)
	return db
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Pass db to handlers, e.g., with dependency injection
	r.GET("/training-days", func(c *gin.Context) {
		GetTrainingDays(c, db)
	})
	r.GET("/training-days/:id", func(c *gin.Context) {
		GetTrainingDayByID(c, db)
	})
	r.POST("/training-days", func(c *gin.Context) {
		CreateTrainingDay(c, db)
	})
	r.PUT("/training-days/:id", func(c *gin.Context) {
		UpdateTrainingDay(c, db)
	})
	r.DELETE("/training-days/:id", func(c *gin.Context) {
		DeleteTrainingDay(c, db)
	})

	return r
}

func SeedTestDatabase(db *gorm.DB) *gorm.DB {
	mesocycle := models.Mesocycle{
		PreparationLevel: 5,
		Comments:         "Test comments",
		Objectives:       "Test objectives",
	}

	if err := db.Create(&mesocycle).Error; err != nil {
		fmt.Println("Error creating Mesocycle:", err)
		return nil
	}
	date := "2024-10-03"
	reps := 50
	weight := 500.0
	trainingDay := models.TrainingDay{
		MesocycleID: mesocycle.ID,
		Date:        &date,
		TotalReps:   &reps,
		TotalWeight: &weight,
	}

	if err := db.Create(&trainingDay).Error; err != nil {
		fmt.Println("Error creating TrainingDay", err)
		return nil
	}

	id := fmt.Sprintf("%d", trainingDay.ID)

	fmt.Printf("TrainingDay: %s inserted\n", id)

	var inserted, err = db.Get(id)
	if err != false {
		fmt.Println("Error creating TrainingDay", err)
		return nil
	}

	fmt.Printf("%+v\n", inserted)

	return db
}

func TestGetTrainingDays(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/training-days", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.TrainingDay
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Greater(t, len(response), 0) // Ensure there are training days in the response
}

func TestGetTrainingDay(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/training-days/1", nil) // Assuming ID 1
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TrainingDay
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "2024-10-03", response.Date) // Check if the correct day is returned
}

// Test AddTrainingDay (POST /training-days)
func TestAddTrainingDay(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	date := "2024-10-10"
	newDay := models.TrainingDay{
		MesocycleID: 1,
		Date:        &date,
	}

	// Marshal the struct into JSON
	jsonBody, err := json.Marshal(newDay)
	if err != nil {
		fmt.Println("Error marshalling the JSON", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/training-days", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.TrainingDay
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "2024-10-10", response.Date) // Check if the new day is created correctly
}

// Test UpdateTrainingDay (PUT /training-days/:day)
func TestUpdateTrainingDay(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	updatedDay := models.TrainingDay{
		ID:        1,
		TotalReps: nil,
	}

	jsonBody, err := json.Marshal(updatedDay)
	if err != nil {
		fmt.Println("Error marshalling the JSON", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/training-days/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TrainingDay
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Ensure TotalReps is not nil (assuming it is a pointer)
	assert.NotNil(t, response.TotalReps)          // Ensure the TotalReps was updated
	assert.Equal(t, "2024-10-03", *response.Date) // // Ensure it updates the correct day
}

// Test DeleteTrainingDay (DELETE /training-days/:day)
func TestDeleteTrainingDay(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/training-days/1", nil) // Assuming ID 1
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Training day deleted successfully", response["message"]) // Ensure the response confirms deletion
}
