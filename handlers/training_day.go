package handlers

import (
	"app/fitness-app-service/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTrainingDays(c *gin.Context, db *gorm.DB) {
	var trainingDays []models.TrainingDay

	// Fetch all training days and preload exercises and sets
	if err := db.Preload("TrainingDayExercises.Sets").Find(&trainingDays).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trainingDays)
}

func GetTrainingDayByID(c *gin.Context, db *gorm.DB) {
	var trainingDay models.TrainingDay
	id := c.Param("id")

	// Find the training day by ID and preload exercises and sets
	if err := db.Preload("TrainingDayExercises.Sets").First(&trainingDay, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training day not found"})
		return
	}

	c.JSON(http.StatusOK, trainingDay)
}

func CreateTrainingDay(c *gin.Context, db *gorm.DB) {
	var trainingDay models.TrainingDay

	// Bind the request JSON to the trainingDay model
	if err := c.ShouldBindJSON(&trainingDay); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the new training day in the database
	if err := db.Create(&trainingDay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trainingDay)
}

func UpdateTrainingDay(c *gin.Context, db *gorm.DB) {
	var trainingDay models.TrainingDay

	fmt.Println("i arrive")
	var input = trainingDay
	id := c.Param("id")
	fmt.Println("Id to update: " + id)
	// Find the existing training day
	if err := db.Preload("TrainingDayExercises.Sets").First(&trainingDay, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training day not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Date != nil {
		trainingDay.Date = input.Date
	}
	if input.TotalReps != nil {
		trainingDay.TotalReps = input.TotalReps
	}
	if input.TotalWeight != nil {
		trainingDay.TotalWeight = input.TotalWeight
	}
	if input.Excercises != nil {
		// Update exercises and sets if provided
		trainingDay.Excercises = input.Excercises
	}

	// Perform the update
	if err := db.Model(&trainingDay).Updates(trainingDay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trainingDay)

	c.JSON(http.StatusOK, trainingDay)
}

func DeleteTrainingDay(c *gin.Context, db *gorm.DB) {
	var trainingDay models.TrainingDay
	id := c.Param("id")

	// Find the training day
	if err := db.First(&trainingDay, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training day not found"})
		return
	}

	// Delete the training day and its related exercises and sets
	if err := db.Delete(&trainingDay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Training day deleted successfully"})
}

// Get all exercises for a specific training day
func GetExercisesByTrainingDayID(c *gin.Context, db *gorm.DB) {
	var exercises []models.Exercise
	trainingDayID := c.Param("trainingDayId")

	if err := db.Where("training_day_id = ?", trainingDayID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercises)
}
