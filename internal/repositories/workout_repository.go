package repositories

import (
	"fmt"
	"strings"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WorkoutRepository struct {
	db *sqlx.DB
}

func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository{
	return &WorkoutRepository{db: db}
}

func (r *WorkoutRepository) SaveWorkout(workout models.Workout) (int, error){
	const op = "internal.repositories.SaveWorkout"
	var workoutID int

	query := `INSERT INTO workouts (user_id, program_id, date, duration, calories, created_at)
		VALUES($1, $2, CURRENT_DATE, $3, $4, NOW()) RETURNING id`
	
	if err := r.db.QueryRow(query, workout.UserID, workout.ProgramID, workout.Duration.Nanoseconds(), workout.Calories).Scan(&workoutID); err != nil{
		return 0, fmt.Errorf("%s: failed to create workout: %w", op, err)
	}

	if err := r.SaveExercisesWorkout(workoutID, workout.Exercises); err != nil{
		return 0, fmt.Errorf("%s: failed to save workout exercises: %w", op, err)
	}

	return workoutID, nil
}

func (r *WorkoutRepository) DeleteWorkout(workoutID int, userID int) (int, error){
	const op = "internal.repositories.DeleteWorkout"
	
	var existID int
	existQuery := `SELECT id FROM workouts WHERE id = $1`
	if err := r.db.Get(&existID, existQuery, workoutID); err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query := `DELETE FROM workouts WHERE id = $1 AND user_id = $2`
	if _, err := r.db.Exec(query, workoutID, userID); err != nil{
		return 0, fmt.Errorf("%s: failed to delete workout: %w", op, err)
	}
	
	return workoutID, nil
}

func (r *WorkoutRepository) DeleteWorkoutExercises(workoutID int) error{
	const op = "internal.repositories.DeleteWorkoutExercises"
	
	query := `DELETE FROM exercises_entry WHERE workout_id = $1`
	if _, err := r.db.Exec(query, workoutID); err != nil{
		return fmt.Errorf("%s: failed to delete workout exercises: %w", op, err)
	}
	
	return nil
}

func (r *WorkoutRepository) GetWorkoutByID(workoutID int, userID int) (*models.Workout, error){
	const op = "internal.repositories.GetWorkoutByID" 
	var workout models.Workout

	query := `SELECT * FROM workouts WHERE id = $1 AND user_id = $2`

	if err := r.db.Get(&workout, query, workoutID, userID); err != nil{
		return nil, fmt.Errorf("%s: failed to get workout: %w", op, err)
	}

	return &workout, nil
}

func (r *WorkoutRepository) GetExercsisesWorkout(workoutID int) ([]models.ExerciseEntry, error){
	const op = "internal.repositories.GetExercsisesWorkout"
	var exercises []models.ExerciseEntry

	query := `SELECT * FROM exercises_entry WHERE workout_id = $1`

	if err := r.db.Select(&exercises, query, workoutID); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return exercises, nil
}

func (r *WorkoutRepository) GetExercisesByID(idSlice []int) ([]models.Exercise, error){
	const op = "internal.repositories.GetExercisesByID"
	var exercises []models.Exercise

	query := `SELECT id, name FROM exercises WHERE id = ANY($1)`
	
	if err := r.db.Select(&exercises, query, pq.Array(idSlice)); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil
}

func (r *WorkoutRepository) SaveExercisesWorkout(workoutID int, exercises []models.ExerciseEntry) error{
	const op = "internal.repositories.SaveExercisesWorkout"

	values := []interface{}{}
	query := `INSERT INTO exercises_entry (workout_id, exercise_id, sets, reps, weight) VALUES `
	placeholderID := 1
	placeholders := []string{}

	for _, ex := range exercises {
		values = append(values, workoutID, ex.ExerciseID, ex.Sets, pq.Array(ex.Reps), pq.Array(ex.Weight))
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

func (r *WorkoutRepository) GetProgramIdByName(programName string) (int, error){
	const op = "internal.repositories.GetProgramIdByName"
	var programID int

	query := `SELECT id FROM programs WHERE name = $1`
	
	if err := r.db.Get(&programID, query, programName); err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return programID, nil
}

func (r *WorkoutRepository) GetExercisesByNames(names []string) ([]models.Exercise, error){
	const op = "internal.repositories.GetExercisesByNames"
	var exercises []models.Exercise

	query := `SELECT id, name FROM exercises WHERE name = ANY($1)`
	if err := r.db.Select(&exercises, query, pq.Array(names)); err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return exercises, nil
}