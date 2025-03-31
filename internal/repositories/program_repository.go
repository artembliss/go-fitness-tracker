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

func (r *ProgramRepository) GetExercisesByID(idSlice []int) ([]models.Exercise, error){
	const op = "internal.repositories.GetExercisesByID"
	var exercises []models.Exercise

	query := `SELECT id, name FROM exercises WHERE id = ANY($1)`
	if err := r.db.Select(&exercises, query, pq.Array(idSlice)); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil
}

func(r *ProgramRepository) GetProgramByID(programID int, userID int) ([]models.Program, error){
	const op = "internal.repositories.GetProgramsByID"
	var programResp []models.Program 
	
	query := `SELECT * FROM programs WHERE id = $1 AND user_id = $2`

	var programsDB []models.ProgramDB 
	if err := r.db.Select(&programsDB, query, programID, userID); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

    for _, programDB := range programsDB {
		var exercises []models.ExerciseProgramDB
        
		if len(programDB.Exercises) > 0 {
			err := json.Unmarshal(programDB.Exercises, &exercises)
            if err != nil {
                return nil, fmt.Errorf("%s: failed to unmarshal exercises: %w", op, err)
            }
			programResp = append(programResp, models.Program{
	        ID:        programDB.ID,
            UserID:    programDB.UserID,
            Name:      programDB.Name,
            Exercises: exercises,
            CreatedAt: programDB.CreatedAt,
			})
        }
    }

	return programResp, nil
}

func (r *ProgramRepository) DeleteProgram(programID int, userID int) (int, error){
	const op = "internal.repositories.DeleteProgram"
	
	var existID int
	queryExist := `SELECT id FROM programs WHERE id = $1`
	if err := r.db.Get(&existID, queryExist, programID); err != nil{
		return 0, fmt.Errorf("%s: program does not exist: %w", op, err)
	}

	var deletedID int
	query := `DELETE FROM programs WHERE id = $1 AND user_id = $2 RETURNING id`
	if err := r.db.Get(&deletedID, query, programID, userID); err != nil{
		return 0, fmt.Errorf("%s: You are not authorized to delete this program: %w", op, err)
	}
	return deletedID, nil
}