package models

import "time"

type Program struct {
	ID        string          `json:"id" db:"id"`
	UserID    string          `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Exercises []ExerciseEntry `json:"exercises" db:"exercises"`
	CreatedAt time.Time       `json:"-" db:"created_at"`
}

type RequestCreateProgram struct {
	Name      string          `json:"name"`
	Exercises []ExerciseEntry `json:"exercises"`
}

type RequestUpdateProgram struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Exercises []ExerciseEntry `json:"exercises"`
}