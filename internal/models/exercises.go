package models

type Exercise struct{
	ID           string       `json:"-" db:"id"`
	Name         string    `json:"-" db:"name"`
	MuscleGroup string    `json:"-" db:"muscle_group"` 
	Description  string    `json:"-" db:"description"`
}

type ExerciseEntry struct {
	ExerciseID string  `json:"exercise_id" db:"exercise_id"`
	Sets       int     `json:"sets" db:"sets"`
	Reps       []int   `json:"reps" db:"reps"`
	Weight     []float64 `json:"weight" db:"weight"`
}