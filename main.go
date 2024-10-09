package main

import (
	"app/fitness-app-service/database"
	"app/fitness-app-service/routes"
)

// func SeedDatabaseMesocycle() {
// 	// Open the JSON file
// 	jsonData, err := ioutil.ReadFile("mesociclo1dia1.json")
// 	if err != nil {
// 		log.Fatal("Failed to read JSON file: ", err)
// 	}

// 	// Parse the JSON data
// 	var mesocycle models.Mesocycle
// 	err = json.Unmarshal(jsonData, &exercises)
// 	if err != nil {
// 		log.Fatal("Failed to unmarshal JSON: ", err)
// 	}

// 	// Insert each exercise into the database
// 	for _, exercise := range exercises {
// 		newExercise := models.Exercise{
// 			Name:         exercise.Name,
// 			MuscleGroup:  exercise.MuscleGroup,
// 			MovementType: exercise.MovementType,
// 			Notes:        exercise.Notes,
// 		}
// 		database.DB.Create(&newExercise)
// 	}

// 	fmt.Println("Database seeded successfully!")
// }

func main() {
	database.InitDatabase()
	// seeder.SeedDatabaseExcercises()
	// Start your API (with Gin router, etc.)
	router := routes.SetupRouter(database.DB)
	router.Run("localhost:8080")
}
