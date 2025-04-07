package repositories

import (
	"fmt"
	"strings"

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

	query := `INSERT INTO programs (user_id, name, created_at)
	        VALUES($1, $2, NOW()) RETURNING id`

	if err := r.db.QueryRow(query, program.UserID, program.Name).Scan(&program.ID); err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.SaveExercisesProgram(program.ID, program.Exercises); err != nil{
		return program.ID, fmt.Errorf("%s: failed to save exercises: %w", op, err)
	}
	
	return program.ID, nil
}

func (r *ProgramRepository) UpdateProgram(program models.Program, programID int) (int, error){
	const op = "internal.repositories.UpdateProgram"

	query := `UPDATE programs SET user_id = $1, name = $2, created_at = NOW()
	          WHERE id = $3 AND user_id = $4
			  RETURNING id`

	if err := r.db.QueryRow(query, program.UserID, program.Name, programID, program.UserID).Scan(&program.ID); err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.SaveExercisesProgram(program.ID, program.Exercises); err != nil{
		return program.ID, fmt.Errorf("%s: failed to save program exercises: %w", op, err)
	}
	
	return program.ID, nil
}

func (r *ProgramRepository) SaveExercisesProgram(programID int, exercises []models.ExerciseProgramDB) error {
	const op = "internal.repositories.SaveExercisesProgram"

	values := []interface{}{}
	query := `INSERT INTO exercises_program (program_id, exercise_id, sets, reps, weight) VALUES `
	placeholderID := 1
	placeholders := []string{}

	for _, ex := range exercises {
		values = append(values, programID, ex.ExerciseID, ex.Sets, ex.Reps, ex.Weight)
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				placeholderID, placeholderID+1, placeholderID+2, placeholderID+3, placeholderID+4),
		)
		placeholderID += 5
	}

	query += strings.Join(placeholders, ", ")

	if _, err := r.db.Exec(query, values...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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

func(r *ProgramRepository) GetProgramByID(programID int, userID int) (*models.Program, error){
	const op = "internal.repositories.GetProgramByID"
	var programDB models.ProgramDB 

	query := `SELECT * FROM programs WHERE id = $1 AND user_id = $2`

	if err := r.db.Get(&programDB, query, programID, userID); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exercises, err := r.GetExercsisesProgram(programID)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	programResp := models.Program{
		ID: programDB.ID,
		UserID: programDB.UserID,
		Name: programDB.Name,
		Exercises: exercises,
		CreatedAt: programDB.CreatedAt,
	}

	return &programResp, nil
}

func (r *ProgramRepository) GetExercsisesProgram(programID int) ([]models.ExerciseProgramDB, error){
	const op = "internal.repositories.GetExercsisesProgram"
	var exercises []models.ExerciseProgramDB

	query := `SELECT * FROM exercises_program WHERE program_id = $1`

	if err := r.db.Select(&exercises, query, programID); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return exercises, nil
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

func (r *ProgramRepository) DeleteExercisesProgram(programID int) (error){
	const op = "internal.repositories.DeleteExercisesProgram"
	
	query := `DELETE FROM exercises_program WHERE program_id = $1`
	if _, err := r.db.Exec(query, programID); err != nil{
		return fmt.Errorf("%s: You are not authorized to delete this program: %w", op, err)
	}
	return nil
}