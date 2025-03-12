package repositories

import (
	"fmt"
	"log"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/storage/postgre"
)

func CheckExercisesExist(s *postgre.Storage) bool {
	var count int

	CheckQuery := `SELECT COUNT(*) FROM exercises`
	err := s.GetDB().Get(&count, CheckQuery)
	if err != nil {
		log.Println("Failed to check storage:", err)
		return false
	}
	return count > 0
}

func SaveExercisesToDB(s *postgre.Storage, exercises []models.ExerciseAPI) (error){
	const op = "internal.handlers.SaveExercisesToDB"

	for _, ex := range exercises{
		_, err := s.GetDB().Exec(`
		INSERT INTO exercises (name, type, muscle_group, equipment, difficulty, instruction)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (name) DO NOTHING`,
		ex.Name, ex.Type, ex.MuscleGroup, ex.Equipment, ex.Difficulty, ex.Instruction)
		if err != nil{
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

