package models

type Exercise struct{
	ID           int    `db:"id"`
	Name         string    `db:"name"`
	Type         string    `db:"type"`
	MuscleGroup  string    `db:"muscle_group"` 
	Equipment    string    `db:"equipment"`
	Difficulty   string    `db:"difficulty"`
	Instruction  string    `db:"instruction"`
}

type ExerciseEntry struct {
	ExerciseID int    `json:"exercise_id" db:"exercise_id"`
	Sets       int       `json:"sets" db:"sets"`
	Reps       []int     `json:"reps" db:"reps"`
	Weight     []float64 `json:"weight" db:"weight"`
}

type ExerciseAPI struct{
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	MuscleGroup  string    `json:"target"` 
	Equipment    string    `json:"equipment"`
	Difficulty   string    `json:"difficulty"`
	Instruction  string    `json:"instructions"`
}
