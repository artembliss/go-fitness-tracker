package models

type Exercise struct{
	ID           int    `json:"-" db:"id"`
	Name         string    `json:"-" db:"name"`
	Type         string    `json:"-" db:"type"`
	MuscleGroup  string    `json:"-" db:"muscle_group"` 
	Equipment    string    `json:"-" db:"equipment"`
	Difficulty   string    `json:"-" db:"difficulty"`
	Instruction  string    `json:"-" db:"instruction"`
}

type ExerciseEntry struct {
	ExerciseID int    `json:"exercise_id" db:"exercise_id"`
	Sets       int       `json:"sets" db:"sets"`
	Reps       []int     `json:"reps" db:"reps"`
	Weight     []float64 `json:"weight" db:"weight"`
}
