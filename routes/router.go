package routes

import (
	"app/fitness-app-service/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Pass db to handlers using dependency injection
	r.GET("/training", func(c *gin.Context) {
		handlers.GetTrainings(c, db)
	})
	r.GET("/training/:id", func(c *gin.Context) {
		handlers.GetTrainingByID(c, db)
	})
	r.POST("/training", func(c *gin.Context) {
		handlers.CreateTraining(c, db)
	})
	r.PUT("/training/:id", func(c *gin.Context) {
		handlers.UpdateTraining(c, db)
	})
	r.DELETE("/training/:id", func(c *gin.Context) {
		handlers.DeleteTraining(c, db)
	})

	// Additional routes...
	return r
}
