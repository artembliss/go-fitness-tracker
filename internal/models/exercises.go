package models

import "github.com/lib/pq"

type Exercise struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Type        string `db:"type"`
	MuscleGroup string `db:"muscle_group"`
	Equipment   string `db:"equipment"`
	Difficulty  string `db:"difficulty"`
	Instruction string `db:"instruction"`
}

type ExerciseEntry struct {
	ID         int             `db:"id"`
	WorkoutID  int             `db:"workout_id"`
	ExerciseID int             `db:"exercise_id"`
	Sets       int             `db:"sets"`
	Reps       pq.Int64Array   `db:"reps" swaggertype:"array,integer"`
	Weight     pq.Float64Array `db:"weight" swaggertype:"array,number"`
}

type ExerciseRequestEntry struct {
	Name   string    `json:"name"`
	Sets   int       `json:"sets"`
	Reps   []int     `json:"reps"`
	Weight []float64 `json:"weight"`
}

type ExerciseRequest struct {
	Name   string  `json:"name"`
	Sets   int     `json:"sets"`
	Reps   int     `json:"reps"`
	Weight float64 `json:"weight"`
}

type ExerciseProgramDB struct {
	ID         int     `db:"id"`
	ProgramID  int     `db:"program_id"`
	ExerciseID int     `db:"exercise_id"`
	Sets       int     `db:"sets"`
	Reps       int     `db:"reps"`
	Weight     float64 `db:"weight"`
}

type ExerciseAPI struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	MuscleGroup string `json:"muscle"`
	Equipment   string `json:"equipment"`
	Difficulty  string `json:"difficulty"`
	Instruction string `json:"instructions"`
}
