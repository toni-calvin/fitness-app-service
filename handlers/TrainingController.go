package handlers

import (
	"app/fitness-app-service/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTrainings(c *gin.Context, db *gorm.DB) {
	var Trainings []models.Training

	// Fetch all training days and preload exercises and sets
	if err := db.Preload("Excercises.Sets").Find(&Trainings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Trainings)
}

func GetTrainingByID(c *gin.Context, db *gorm.DB) {
	var Training models.Training
	id := c.Param("id")

	fmt.Println("Id to update: " + id)

	// Find the training day by ID and preload exercises and sets
	if err := db.Preload("Excercises.Sets").First(&Training, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training day not found"})
		return
	}

	fmt.Printf("Retrieved Training: %+v\n", Training)

	c.JSON(http.StatusOK, Training)
}

func CreateTraining(c *gin.Context, db *gorm.DB) {
	var Training models.Training

	// Bind the request JSON to the Training model
	if err := c.ShouldBindJSON(&Training); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the new training day in the database
	if err := db.Create(&Training).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Training)
}

func UpdateTraining(c *gin.Context, db *gorm.DB) {
	var Training models.Training

	id := c.Param("id")
	fmt.Println("Id to update: " + id)
	fmt.Printf("%+v\n", Training)

	// Find the existing training day
	if err := db.Preload("Excercises.Sets").First(&Training, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training day not found"})
		return
	}

	if err := c.ShouldBindJSON(&Training); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform the update
	if err := db.Model(&Training).Updates(Training).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Training)
}

func DeleteTraining(c *gin.Context, db *gorm.DB) {
	var Training models.Training
	id := c.Param("id")

	// Find the training day
	if err := db.First(&Training, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training day not found"})
		return
	}

	// Delete the training day and its related exercises and sets
	if err := db.Delete(&Training).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Training day deleted successfully"})
}

// Get all exercises for a specific training day
func GetExercisesByTrainingID(c *gin.Context, db *gorm.DB) {
	var exercises []models.Exercise
	TrainingID := c.Param("TrainingId")

	if err := db.Where("training_day_id = ?", TrainingID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercises)
}
