package models

import "time"

type Program struct {
	ID        int          `json:"id" db:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Exercises []ExerciseEntry `json:"exercises" db:"exercises"`
	CreatedAt time.Time       `json:"-" db:"created_at"`
}

type RequestCreateProgram struct {
	Name      string          `json:"name"`
	Exercises []ExerciseEntry `json:"exercises"`
}

type RequestUpdateProgram struct {
	ID        int          `json:"id"`
	Name      string          `json:"name"`
	Exercises []ExerciseEntry `json:"exercises"`
}