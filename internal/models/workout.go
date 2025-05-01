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
	Duration  time.Duration   `json:"duration" db:"duration" swaggertype:"integer"`
	Calories  float64         `json:"calories" db:"calories"`
	CreatedAt time.Time       `json:"-" db:"created_at"`
}

type RequestCreateWorkout struct {
	ProgramName string                 `json:"program_name"`
	Exercises   []ExerciseRequestEntry `json:"exercises"`
	Duration    string                 `json:"duration" binding:"required"`
	Calories    float64                `json:"calories"`
}

type RequestGetWorkout struct {
	ID        int                    `json:"id"`
	UserID    int                    `json:"user_id"`
	ProgramID int                    `json:"program_id"`
	Date      time.Time              `json:"date"`
	Exercises []ExerciseRequestEntry `json:"exercises"`
	Duration  string                 `json:"duration"`
	Calories  float64                `json:"calories"`
	CreatedAt time.Time              `json:"-"`
}