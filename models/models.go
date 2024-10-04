package models

import (
	"errors"

	"gorm.io/gorm"
)

type Exercise struct {
	gorm.Model
	Name         string `json:"name" gorm:"column:name"`
	MuscleGroup  string `json:"muscle_group" gorm:"column:muscle_group"`
	MovementType string `json:"movement_type" gorm:"column:movement_type"`
	Notes        string `json:"notes" gorm:"column:notes"`
}

type Mesocycle struct {
	gorm.Model
	PreparationLevel int           `json:"preparation_level" gorm:"column:preparation_level"`
	Comments         string        `json:"comments" gorm:"column:comments"`
	Objectives       string        `json:"objectives" gorm:"column:objectives"`
	TrainingDays     []TrainingDay `json:"training_days" gorm:"foreignKey:MesocycleID"`
}

// gorm hook
func (e *Mesocycle) BeforeSave(tx *gorm.DB) (err error) {
	// Ensure the preparation level is between 1 and 10
	if e.PreparationLevel < 1 || e.PreparationLevel > 10 {
		return errors.New("preparation level must be between 1 and 10")
	}
	return
}

type TrainingDay struct {
	ID          uint                  `json:"id" gorm:"primaryKey"`
	MesocycleID uint                  `json:"mesocycle_id" gorm:"index;foreignKey:MesocycleID"`
	Date        *string               `json:"date" gorm:"column:date"`
	TotalReps   *int                  `json:"total_reps" gorm:"column:total_reps"`
	TotalWeight *float64              `json:"total_weight" gorm:"column:total_weight"`
	Excercises  []TrainingDayExercise `json:"exercises" gorm:"foreignKey:TrainingDayID"`
}

type TrainingDayExercise struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	TrainingDayID uint    `json:"training_day_id" gorm:"index"`
	ExerciseID    uint    `json:"exercise_id" gorm:"index;foreignKey:ExerciseID"`
	TotalReps     int     `json:"total_reps" gorm:"column:total_reps"`
	TotalWeight   float64 `json:"total_weight" gorm:"column:total_weight"`
	Notes         string  `json:"notes" gorm:"column:notes"`
	Sets          []Set   `json:"sets" gorm:"foreignKey:TrainingDayExerciseID"`
}

type Set struct {
	ID                    uint    `json:"id" gorm:"primaryKey"`
	TrainingDayExerciseID uint    `json:"training_day_exercise_id" gorm:"index;foreignKey:TrainingDayExerciseID"`
	SetNumber             int     `json:"set_number" gorm:"column:set_number"`
	Reps                  int     `json:"reps" gorm:"column:reps"`
	RIR                   int     `json:"rir" gorm:"column:rir"`
	Weight                float64 `json:"weight" gorm:"column:weight"`
	RestTime              int     `json:"rest_time" gorm:"column:rest_time"`
}
