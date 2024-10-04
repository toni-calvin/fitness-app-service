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

func CleanTestDatabase(db *gorm.DB) {
	db.Delete(models.Training{})
	db.Delete(&models.Set{})
	db.Delete(&models.TrainingExercise{})
	db.Delete(&models.Training{})
	db.Delete(&models.Mesocycle{})
}

func SetupTestDatabase() *gorm.DB {
	fmt.Println("Setting up db")
	db, _ := gorm.Open(sqlite.Open("../testfitnessapp.db"), &gorm.Config{})
	CleanTestDatabase(db)
	db.AutoMigrate(&models.Mesocycle{}, &models.TrainingExercise{}, &models.Training{}, &models.Exercise{}, &models.Set{})
	db = SeedTestDatabase(db)
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
	}

	if err := db.Create(&mesocycle).Error; err != nil {
		fmt.Println("Error creating Mesocycle:", err)
		return nil
	}
	Training := models.Training{
		MesocycleID: mesocycle.ID,
		Date:        "2024-10-03",
		TotalReps:   50,
		TotalWeight: 500.0,
	}

	if err := db.Create(&Training).Error; err != nil {
		fmt.Println("Error creating Training", err)
		return nil
	}

	id := fmt.Sprintf("%d", Training.ID)

	fmt.Printf("Training: %s inserted\n", id)
	var inserted, err = db.Get(id)
	if err != false {
		fmt.Println("Error retrieving Training", err)
		return nil
	}
	fmt.Printf("%+v\n", inserted)

	return db
}

func TestGetTrainings(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

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
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/training/1", nil) // Assuming ID 1
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Training
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "2024-10-03", response.Date) // Check if the correct day is returned
}

// Test AddTraining (POST /training)
func TestAddTraining(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	newDay := models.Training{
		MesocycleID: 1,
		Date:        "2024-10-10",
	}

	// Marshal the struct into JSON
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
	assert.Equal(t, "2024-10-10", response.Date) // Check if the new day is created correctly
}

// Test UpdateTraining (PUT /training/:day)
func TestUpdateTraining(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router
	updatedDay := models.Training{
		ID:        1,
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

	// Ensure TotalReps is not nil (assuming it is a pointer)
	assert.Zero(t, response.TotalReps)           // Ensure the TotalReps was updated
	assert.Equal(t, "2024-10-03", response.Date) // // Ensure it updates the correct day
}

// Test DeleteTraining (DELETE /training/:day)
func TestDeleteTraining(t *testing.T) {
	db := SetupTestDatabase() // Initialize the test database
	router := SetupRouter(db) // Pass the db into the router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/training/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Training day deleted successfully", response["message"]) // Ensure the response confirms deletion
}
