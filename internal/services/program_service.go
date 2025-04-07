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

func (s *ProgramService) CreateProgram(program models.Program) (int, error){
	const op = "internal.servises.SaveProgram"

	id, err := s.ProgramRepo.SaveProgram(program)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *ProgramService) UpdateProgram(program models.Program, programID int) (int, error){
	const op = "internal.servises.UpdateProgram"

	if err := s.ProgramRepo.DeleteExercisesProgram(programID); err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)	
	}

	id, err := s.ProgramRepo.UpdateProgram(program, programID)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *ProgramService) GetNameToID(exercises []models.ExerciseRequest) (map[string]int, error){
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

func (s *ProgramService) MapToDBExercises(regEx []models.ExerciseRequest, nameToDB map[string]int) ([]models.ExerciseProgramDB, []string){
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

func (s *ProgramService) GetProgram(programID int, userID int) (*models.RequestGetProgram, error){
	const op = "internal.servises.GetPrograms"

	var program *models.RequestGetProgram

	programDB, err := s.ProgramRepo.GetProgramByID(programID, userID)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to get programs by id: %w", op, err)
	}

	program, err = s.BuildResponseExercises(*programDB)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return program, nil
}

func (s *ProgramService) GetIdToName(exercises []models.ExerciseProgramDB) (map[int]string, error){
	const op = "internal.servises.GetIdToName"
	idSlice := make([]int, 0, len(exercises))

	for _, exercise := range exercises{
		idSlice = append(idSlice, exercise.ExerciseID)
	}

	found, err := s.ProgramRepo.GetExercisesByID(idSlice)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to get exercises by idSlice: %w", op, err)
	}

	exerciseMap := make(map[int]string, len(found))
	
	for _, e := range found {
        exerciseMap[e.ID] = e.Name
    }
	
	return exerciseMap, nil
}

func (s *ProgramService) MapToResponseExercises(dbEx []models.ExerciseProgramDB, idToName map[int]string) (*[]models.ExerciseRequest, []int){
	var result []models.ExerciseRequest
	var notFound []int

	for _, ex := range dbEx{
		name, ok := idToName[ex.ExerciseID]
        if !ok {
            notFound = append(notFound, ex.ExerciseID)
            continue
        }
		result = append(result, models.ExerciseRequest{
            Name: name,
            Sets:       ex.Sets,
            Reps:       ex.Reps,
            Weight:     ex.Weight,
        })
	}
	
	return &result, notFound
}

func (s *ProgramService) BuildResponseExercises(programDB models.Program) (*models.RequestGetProgram, error){
	const op = "internal.servises.BuildResponseExercises"

	idToName, err := s.GetIdToName(programDB.Exercises)
	if err != nil{
	 return nil, fmt.Errorf("%s: %w", op, err)
	}

	exercisesResp, notFound := s.MapToResponseExercises(programDB.Exercises, idToName)
	if len(notFound) > 0{
		return nil, fmt.Errorf("%s: some exercises not found: %v", op, notFound)
	}

	programsResp := models.RequestGetProgram{
			ID: programDB.ID,
			UserID: programDB.UserID,
			Name: programDB.Name,
			Exercises: *exercisesResp,
			CreatedAt: programDB.CreatedAt,
		}
	
	return &programsResp, nil
}

func (s *ProgramService) DeleteProgram(programID int, userID int) (int, error){
	const op = "internal.servises.DeleteProgram"
	deletedID, err := s.ProgramRepo.DeleteProgram(programID, userID)
	
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if err := s.ProgramRepo.DeleteExercisesProgram(programID); err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)	
	}
	
	return deletedID, nil
}