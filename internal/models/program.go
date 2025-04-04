package models

import (
	"time"
)

type Program struct {
	ID        int          `json:"id" db:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Exercises []ExerciseProgramDB `json:"exercises" db:"exercises"`
	CreatedAt time.Time       `json:"-" db:"created_at"`
}

type ProgramDB struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	Name      string       `db:"name"`
	CreatedAt time.Time    `db:"created_at"`
}

type RequestGetProgram struct {
	ID        int                     `json:"id"`
	UserID    int                     `json:"user_id"`
	Name      string                  `json:"name"`
	Exercises []ExerciseRequest       `json:"exercises"`
	CreatedAt time.Time               `json:"-"`
}

type RequestCreateProgram struct {
	Name      string          `json:"name"`
	Exercises []ExerciseRequest `json:"exercises"`
}

type RequestUpdateProgram struct {
	ID        int          `json:"id"`
	Name      string          `json:"name"`
	Exercises []ExerciseRequest `json:"exercises"`
}