package routes

import (
	"app/fitness-app-service/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Pass db to handlers using dependency injection
	r.GET("/training-days", func(c *gin.Context) {
		handlers.GetTrainingDays(c, db)
	})
	r.GET("/training-days/:id", func(c *gin.Context) {
		handlers.GetTrainingDayByID(c, db)
	})
	r.POST("/training-days", func(c *gin.Context) {
		handlers.CreateTrainingDay(c, db)
	})
	r.PUT("/training-days/:id", func(c *gin.Context) {
		handlers.UpdateTrainingDay(c, db)
	})
	r.DELETE("/training-days/:id", func(c *gin.Context) {
		handlers.DeleteTrainingDay(c, db)
	})

	// Additional routes...
	return r
}
