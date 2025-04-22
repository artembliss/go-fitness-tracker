package services

import (
	"fmt"
	"time"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/lib/pq"
)

type WorkoutService struct {
	WorkoutRepo *repositories.WorkoutRepository
}

func NewWorkoutService(repo *repositories.WorkoutRepository) *WorkoutService{
	return &WorkoutService{WorkoutRepo: repo}
}

func (s *WorkoutService) CreateWorkout(userID int, workoutCreate models.RequestCreateWorkout) (int, error){
	const op = "internal.servises.CreateWorkout"

	var workout models.Workout
	
	programID, err := s.WorkoutRepo.GetProgramIdByName(workoutCreate.ProgramName)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	nameToID, err := s.GetNameToID(workoutCreate.Exercises)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	exercisesEntryToSave, notFound := s.MapToDBExercisesEntry(workoutCreate.Exercises, nameToID)
	if len(notFound) > 0 {
		return 0, fmt.Errorf("%s: some exercises not found: %w", op, err)
	}

	duration, err := time.ParseDuration(workoutCreate.Duration)
	if err != nil {
		return 0, fmt.Errorf("%s: Invalid duration format: %w", op, err)
	}

	workout = models.Workout{
		UserID: userID,
		ProgramID: programID,
		Exercises: exercisesEntryToSave,
		Duration: duration,
		Calories: workoutCreate.Calories,
	}

	workoutID, err := s.WorkoutRepo.SaveWorkout(workout)
	if err != nil{
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return workoutID, nil
}

func (s *WorkoutService) GetWorkout(workoutID int, userID int) (*models.RequestGetWorkout, error){
	const op = "internal.servises.workout_service.GetWorkout"

	workoutDB, err := s.WorkoutRepo.GetWorkoutByID(workoutID, userID) 
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	workoutDB.Exercises, err = s.WorkoutRepo.GetExercsisesWorkout(workoutID)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	workout, err := s.BuildResponseWorkout(*workoutDB)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return workout, nil
}

func (s *WorkoutService) BuildResponseWorkout(workoutDB models.Workout) (*models.RequestGetWorkout, error){
	const op = "internal.servises.BuildResponseWorkout"
	IdToName, err := s.GetIdToName(workoutDB.Exercises)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to build response: %w", op, err)
	}
	fmt.Println(IdToName)
	exercises, notFound := s.MapToResponseExercises(workoutDB.Exercises, IdToName)
	if len(notFound) > 0{
		return nil, fmt.Errorf("%s: some exercises not found: %v", op, notFound)
	}
	workout := models.RequestGetWorkout{
		ID: workoutDB.ID,
		UserID: workoutDB.UserID,
		ProgramID: workoutDB.ProgramID,
		Date: workoutDB.Date,
		Exercises: exercises,
		Duration: workoutDB.Duration.String(),
		Calories: workoutDB.Calories,
		CreatedAt: workoutDB.CreatedAt,
	}
	return &workout, nil
}

func (s *WorkoutService) MapToResponseExercises(dbEx []models.ExerciseEntry, idToName map[int]string) ([]models.ExerciseRequestEntry, []int) {
    var result []models.ExerciseRequestEntry
    var notFound []int

    for _, ex := range dbEx {
        name, ok := idToName[ex.ExerciseID]
        if !ok {
            notFound = append(notFound, ex.ExerciseID)
            continue
        }

        reps := make([]int, len(ex.Reps))
        for i, v := range ex.Reps {
            reps[i] = int(v) 
        }

        weight := make([]float64, len(ex.Weight))
        copy(weight, ex.Weight)

        result = append(result, models.ExerciseRequestEntry{
            Name:       name,
            Sets:       ex.Sets,
            Reps:       reps,
            Weight:     weight,
        })
    }

    return result, notFound
}



func (s *WorkoutService) GetIdToName(exercises []models.ExerciseEntry) (map[int]string, error){
	op := "internal.servises.workout_service.GetIdToName"
	idSlice := make([]int, len(exercises))

	for _, exercise := range exercises{
		idSlice = append(idSlice, exercise.ExerciseID)
	}

	found, err := s.WorkoutRepo.GetExercisesByID(idSlice) 
	if err != nil{
		return nil, fmt.Errorf("%s: failed to get exercises by idSlice: %w", op, err)
	}

	exerciseMap := make(map[int]string, len(found))

	for _, ex := range found{
		exerciseMap[ex.ID] = ex.Name
	}

	return exerciseMap, nil
}

func (s *WorkoutService) GetNameToID(exercises []models.ExerciseRequestEntry) (map[string]int, error){
	const op = "internal.servises.workout_service.GetNameToID"
	names := make([]string, len(exercises))

	for _, exercise := range exercises{
		names = append(names, exercise.Name)
	}

	found, err := s.WorkoutRepo.GetExercisesByNames(names)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to get exercises by names: %w", op, err)
	}

	exerciseMap := make(map[string]int, len(found))
	
	for _, e := range found {
        exerciseMap[e.Name] = e.ID
    }
	
	return exerciseMap, nil
}

func (s *WorkoutService) MapToDBExercisesEntry(regEx []models.ExerciseRequestEntry, nameToDB map[string]int) ([]models.ExerciseEntry, []string) {
    var result []models.ExerciseEntry
    var notFound []string

    for _, ex := range regEx {
        id, ok := nameToDB[ex.Name]
        if !ok {
            notFound = append(notFound, ex.Name)
            continue
        }

        reps := make(pq.Int64Array, len(ex.Reps))
        for i, v := range ex.Reps {
            reps[i] = int64(v) 
        }

        weight := make(pq.Float64Array, len(ex.Weight))
        copy(weight, ex.Weight)  

        result = append(result, models.ExerciseEntry{
            ExerciseID: id,
            Sets:       ex.Sets,
            Reps:       reps,
            Weight:     weight,
        })
    }

    return result, notFound
}