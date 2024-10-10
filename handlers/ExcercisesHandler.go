package handlers

import (
	"app/fitness-app-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetExcercises(c *gin.Context, db *gorm.DB) {
	var Excercises []models.Exercise
	if err := db.Find(&Excercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Excercises)
}
