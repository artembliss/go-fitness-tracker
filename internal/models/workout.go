package models

import (
	"time"
)

type Workout struct {
	ID        int          `json:"id" db:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	ProgramID int          `json:"program_id" db:"program_id"`
	Date      time.Time       `json:"date" db:"date"`
	Exercises []ExerciseEntry `json:"exercises" db:"exercises"`
	Duration  int             `json:"duration" db:"duration"`
	Calories  float64         `json:"calories" db:"calories"`
	CreatedAt time.Time       `json:"-" db:"created_at"`
}

type RequestCreateWorkout struct {
	ProgramName string          `json:"program_name"`
	Exercises   []ExerciseEntry `json:"exercises"`
}

type RequestUpdateWorkout struct {
	ID         int             `json:"id"`
	ProgramName string         `json:"program_name"`
	Exercises  []ExerciseEntry `json:"exercises"`
}