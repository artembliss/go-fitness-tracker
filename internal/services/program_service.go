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

func (s *ProgramService) GetPrograms(userId int) (*[]models.RequestGetProgram, error){
	const op = "internal.servises.GetPrograms"

	var programs []models.RequestGetProgram 

	programsDB, err := s.ProgramRepo.GetProgramsByID(userId)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to get programs by id: %w", op, err)
	}

	programs, err = s.BuildResponseExercises(programsDB)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &programs, nil
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

func (s *ProgramService) MapToResponseExercises(dbEx []models.ExerciseProgramDB, idToName map[int]string) (*[]models.ExerciseRequestCreate, []int){
	var result []models.ExerciseRequestCreate
	var notFound []int

	for _, ex := range dbEx{
		name, ok := idToName[ex.ExerciseID]
        if !ok {
            notFound = append(notFound, ex.ExerciseID)
            continue
        }
		result = append(result, models.ExerciseRequestCreate{
            Name: name,
            Sets:       ex.Sets,
            Reps:       ex.Reps,
            Weight:     ex.Weight,
        })
	}
	
	return &result, notFound
}

func (s *ProgramService) BuildResponseExercises(programsDB []models.Program) ([]models.RequestGetProgram, error){
	const op = "internal.servises.BuildResponseExercises"
	var programsResp []models.RequestGetProgram

	for _, programDB := range programsDB{
		idToName, err := s.GetIdToName(programDB.Exercises)
		if err != nil{
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		exercisesResp, notFound := s.MapToResponseExercises(programDB.Exercises, idToName)
		if len(notFound) > 0{
			return nil, fmt.Errorf("%s: some exercises not found: %v", op, notFound)
		}

		programsResp = append(programsResp, models.RequestGetProgram{
			ID: programDB.ID,
			UserID: programDB.UserID,
			Name: programDB.Name,
			Exercises: *exercisesResp,
			CreatedAt: programDB.CreatedAt,
		})
	}
	return programsResp, nil
}