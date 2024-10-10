package routes

import (
	"app/fitness-app-service/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"}, // Add the origin of your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

	r.GET("/exercises", func(c *gin.Context) {
		handlers.GetExcercises(c, db)
	})

	r.POST("/mesocycles", func(c *gin.Context) {
		handlers.CreateMesocycle(c, db)
	})

	return r
}
