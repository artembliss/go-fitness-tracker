package repositories

import (
	"fmt"
	"log"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type ExerciseRepository struct {
	db *sqlx.DB
}

func NewExerciseRepository(db *sqlx.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) CheckExercisesExist() bool {
	var count int

	CheckQuery := `SELECT COUNT(*) FROM exercises`
	err := r.db.Get(&count, CheckQuery)
	if err != nil {
		log.Println("Failed to check storage:", err)
		return false
	}
	return count > 0
}

func (r *ExerciseRepository) SaveExercisesToDB(exercises []models.ExerciseAPI) (error){
	const op = "internal.handlers.SaveExercisesToDB"

	for _, ex := range exercises{
		_, err := r.db.Exec(`
		INSERT INTO exercises (name, type, muscle_group, equipment, difficulty, instruction)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (name) DO NOTHING`,
		ex.Name, ex.Type, ex.MuscleGroup, ex.Equipment, ex.Difficulty, ex.Instruction)
		if err != nil{
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (r *ExerciseRepository) GetAllExercises() ([]models.Exercise, error){
	op := "internal.handlers.GetAllExercises"

	var exercises []models.Exercise

	getAllExercisesQuery := `SELECT id, name, type, muscle_group, equipment, difficulty, instruction FROM exercises`
	if err := r.db.Select(&exercises, getAllExercisesQuery); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(exercises) == 0{
		return nil, fmt.Errorf("%s: storage is empty", op)
	}
	return exercises, nil
}
