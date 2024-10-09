package models

import (
	"errors"

	"gorm.io/gorm"
)

type Exercise struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"Name"`
	MuscleGroup  string `json:"MuscleGroup"`
	MovementType string `json:"MovementType"`
	Notes        string `json:"Notes"`
}

type Mesocycle struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	PreparationLevel int        `json:"preparation_level" gorm:"column:preparation_level"`
	Comments         string     `json:"comments" gorm:"column:comments"`
	Objectives       string     `json:"objectives" gorm:"column:objectives"`
	Trainings        []Training `json:"training_days" gorm:"foreignKey:MesocycleID"`
}

// gorm hook
func (e *Mesocycle) BeforeSave(tx *gorm.DB) (err error) {
	// Ensure the preparation level is between 1 and 10
	if e.PreparationLevel < 1 || e.PreparationLevel > 10 {
		return errors.New("preparation level must be between 1 and 10")
	}
	return
}

type Training struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	MesocycleID uint               `json:"mesocycle_id" gorm:"index;foreignKey:MesocycleID"`
	Date        string             `json:"date" gorm:"column:date"`
	TotalReps   int                `json:"total_reps" gorm:"column:total_reps"`
	TotalWeight float64            `json:"total_weight" gorm:"column:total_weight"`
	Excercises  []TrainingExercise `json:"exercises" gorm:"foreignKey:TrainingID"`
}

type TrainingExercise struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	TrainingID  uint     `json:"training_day_id" gorm:"index"`
	ExerciseID  uint     `json:"exercise_id" gorm:"index"`
	Exercise    Exercise `json:"exercise" gorm:"foreignKey:ExerciseID"`
	TotalReps   int      `json:"total_reps"`
	TotalWeight float64  `json:"total_weight"`
	Notes       string   `json:"notes"`
	Sets        []Set    `json:"sets" gorm:"foreignKey:TrainingExerciseID"`
}

type Set struct {
	ID                 uint    `json:"id" gorm:"primaryKey"`
	TrainingExerciseID uint    `json:"training_day_exercise_id" gorm:"index;foreignKey:TrainingExerciseID"`
	SetNumber          int     `json:"set_number" gorm:"column:set_number"`
	Reps               int     `json:"reps" gorm:"column:reps"`
	RIR                int     `json:"rir" gorm:"column:rir"`
	Weight             float64 `json:"weight" gorm:"column:weight"`
	RestTime           int     `json:"rest_time" gorm:"column:rest_time"`
}
