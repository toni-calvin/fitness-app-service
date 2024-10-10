package models

type Set struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Reps       int    `json:"reps"`
	Weight     string `json:"weight"`
	RestTime   int    `json:"restTime"`
	RIR        int    `json:"rir"`
	ExerciseID uint   `json:"exerciseId"`
}

type Exercise struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Sets         []Set  `json:"sets" gorm:"foreignKey:ExerciseID"`
	MuscleGroup  string `json:"muscleGroup"`
	MovementType string `json:"movementType"`
	Notes        string `json:"notes"`
	TrainingID   uint   `json:"trainingId"`
}

type Training struct {
	ID                int        `json:"id" gorm:"primaryKey"`
	TargetMuscleGroup string     `json:"targetMuscleGroup"`
	Exercises         []Exercise `json:"exercises" gorm:"foreignKey:TrainingID"`
	PreparationLevel  int        `json:"preparationLevel"`
	Comments          string     `json:"comments"`
	MicrocycleID      uint       `json:"microcycleId"`
}

type Microcycle struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	StartDate   string     `json:"startDate"`
	EndDate     string     `json:"endDate"`
	Trainings   []Training `json:"trainings" gorm:"foreignKey:MicrocycleID"`
	MesocycleID uint       `json:"mesocycleId"`
}

type Mesocycle struct {
	ID               int          `json:"id" gorm:"primaryKey"`
	StartDate        string       `json:"startDate"`
	EndDate          string       `json:"endDate"`
	Microcycles      []Microcycle `json:"microcycles" gorm:"foreignKey:MesocycleID"`
	PreparationLevel int          `json:"preparationLevel"`
	Comments         string       `json:"comments"`
	Objectives       string       `json:"objectives"`
}
