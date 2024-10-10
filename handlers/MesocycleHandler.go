package handlers

import (
	"app/fitness-app-service/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMesocycle(c *gin.Context, db *gorm.DB) {
	var mesocycleForm models.CreateMesocycleForm
	fmt.Println("Creating Mesocycle: %+v\n", mesocycleForm)
	if err := c.ShouldBindJSON(&mesocycleForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new Mesocycle object
	newMesocycle := models.Mesocycle{
		StartDate:  mesocycleForm.StartDate,
		Objectives: mesocycleForm.Objectives,
	}
	numberMicrocycles, err := strconv.Atoi(mesocycleForm.NumberMicrocycles)
	if err != nil {
		fmt.Println("Error converting number of microcycles to int")
	}
	for i := 0; i < numberMicrocycles; i++ {
		microcycle := models.Microcycle{}
		newMesocycle.Microcycles = append(newMesocycle.Microcycles, microcycle)
	}

	if err := db.Create(&newMesocycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newMesocycle)
}
