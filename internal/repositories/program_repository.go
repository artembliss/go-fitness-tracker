package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProgramRepository struct {
	db *sqlx.DB
}

func NewProgramRepository(db *sqlx.DB) *ProgramRepository{
	return &ProgramRepository{db: db}
} 

func (r *ProgramRepository) SaveProgram(program models.Program) (int, error){
	const op = "internal.repositories.SaveProgram"

	jsonExercises, err := json.Marshal(program.Exercises)
	if err != nil {
		return 0,fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO programs (user_id, name, exercises, created_at)
	        VALUES($1, $2, $3, NOW()) RETURNING id`

	err = r.db.QueryRow(query, program.UserID, program.Name, jsonExercises).Scan(&program.ID)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)

	}
	
	return program.ID, nil
}


func (r *ProgramRepository) GetExercisesByNames(names []string) ([]models.Exercise, error){
	const op = "internal.repositories.GetExercisesByNames"
	var exercises []models.Exercise

	query := `SELECT id, name FROM exercises WHERE name = ANY($1)`
	if err := r.db.Select(&exercises, query, pq.Array(names)); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil
}