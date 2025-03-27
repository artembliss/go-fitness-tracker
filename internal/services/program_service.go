package services

import (
	"fmt"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
)

type ProgramService struct {
	ProgramRepo *repositories.ProgramRepository
}

func NewProgramService(repo *repositories.ProgramRepository) *ProgramService{
	return &ProgramService{ProgramRepo: repo}
}

func (s *ProgramService) SaveProgram(program models.Program) (int, error){
	const op = "internal.servises.SaveProgram"
	
	id, err := s.ProgramRepo.SaveProgram(program)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *ProgramService) GetNameToID(exercises []models.ExerciseRequestCreate) (map[string]int, error){
	const op = "internal.servises.GetNameToID"
	names := make([]string, len(exercises))

	for _, exercise := range exercises{
		names = append(names, exercise.Name)
	}

	found, err := s.ProgramRepo.GetExercisesByNames(names)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to get exercises by names: %w", op, err)
	}

	exerciseMap := make(map[string]int, len(found))
	
	for _, e := range found {
        exerciseMap[e.Name] = e.ID
    }
	
	return exerciseMap, nil
}

func (s *ProgramService) MapToDBExercises(regEx []models.ExerciseRequestCreate, nameToDB map[string]int) ([]models.ExerciseProgramDB, []string){
	var result []models.ExerciseProgramDB
	var notFound []string

	for _, ex := range regEx{
		id, ok := nameToDB[ex.Name]
        if !ok {
            notFound = append(notFound, ex.Name)
            continue
        }
		result = append(result, models.ExerciseProgramDB{
            ExerciseID: id,
            Sets:       ex.Sets,
            Reps:       ex.Reps,
            Weight:     ex.Weight,
        })
	}
	
	return result, notFound
}