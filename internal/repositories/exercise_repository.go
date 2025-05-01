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
	const op = "internal.repositories.SaveExercisesToDB"

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
	op := "internal.repositories.GetAllExercises"

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

func (r *ExerciseRepository) GetExercisesByID(id int) (models.Exercise, error){
	op := "internal.repositories.GetExercisesByID"

	var exercise models.Exercise

	getAllExercisesQuery := `SELECT id, name, type, muscle_group, equipment, difficulty, instruction FROM exercises
	WHERE ID = $1`
	if err := r.db.Get(&exercise, getAllExercisesQuery, id); err != nil{
		return models.Exercise{}, fmt.Errorf("%s: %w", op, err)
	}
	if exercise == (models.Exercise{}){
		return models.Exercise{}, fmt.Errorf("%s: Invalid id", op)
	}

	return exercise, nil
}

func (r *ExerciseRepository) GetExercisesByName(id string) (models.Exercise, error){
	op := "internal.repositories.GetExercisesByName"

	var exercise models.Exercise

	getExercisesByNameQuery := `SELECT id, name, type, muscle_group, equipment, difficulty, instruction FROM exercises
	WHERE name = $1`
	if err := r.db.Get(&exercise, getExercisesByNameQuery, id); err != nil{
		return models.Exercise{}, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercise, nil
}

func (r *ExerciseRepository) GetExercisesByType(typeEx string) ([]models.Exercise, error){
	op := "internal.repositories.GetExercisesByType"

	var exercises []models.Exercise

	getExercisesByNameQuery := `SELECT id, name, type, muscle_group, equipment, difficulty, instruction FROM exercises
	WHERE type = $1`
	if err := r.db.Select(&exercises, getExercisesByNameQuery, typeEx); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil	
}

func (r *ExerciseRepository) GetExercisesByMuscleGroup(muscleGroup string) ([]models.Exercise, error){
	op := "internal.repositories.GetExercisesByMuscleGroup"

	var exercises []models.Exercise

	getExercisesByNameQuery := `SELECT id, name, type, muscle_group, equipment, difficulty, instruction FROM exercises
	WHERE muscle_group = $1`
	if err := r.db.Select(&exercises, getExercisesByNameQuery, muscleGroup); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil	
}

func (r *ExerciseRepository) GetExercisesByDifficulty(difficulty string) ([]models.Exercise, error){
	op := "internal.repositories.GetExercisesByDifficulty"

	var exercises []models.Exercise

	getExercisesByNameQuery := `SELECT id, name, type, muscle_group, equipment, difficulty, instruction FROM exercises
	WHERE Difficulty = $1`
	if err := r.db.Select(&exercises, getExercisesByNameQuery, difficulty); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil	
}