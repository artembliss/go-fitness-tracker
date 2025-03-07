package models

import "time"

type Workout struct {
	ID        string          `json:"id" db:"id"`
	UserID    string          `json:"user_id" db:"user_id"`
	ProgramID string          `json:"program_id" db:"program_id"`
	Exercises []ExerciseEntry `json:"exercises" db:"exercises"`
	CreatedAt time.Time       `json:"-" db:"created_at"`
	Duration  int             `json:"duration" db:"duration"`
	Calories  float64         `json:"calories" db:"calories"`
}

type RequestCreateWorkout struct {
	ProgramName string          `json:"program_name"`
	Exercises   []ExerciseEntry `json:"exercises"`
}

type RequestUpdateWorkout struct {
	ID         string             `json:"id"`
	ProgramName string         `json:"program_name"`
	Exercises  []ExerciseEntry `json:"exercises"`
}