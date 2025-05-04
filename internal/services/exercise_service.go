package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/redis/go-redis/v9"
)

type ExerciseService struct{
	ExerciseRepo *repositories.ExerciseRepository
	Cache        *redis.Client
}

func NewExerciseService(repo *repositories.ExerciseRepository, cache *redis.Client) *ExerciseService {
	return &ExerciseService{
		ExerciseRepo: repo,
	    Cache: cache,
	}
}

type ServiceFunc func(param ...interface{}) (interface{}, error)

func (s *ExerciseService) FetchExercisesByMuscle(muscle string) ([]models.ExerciseAPI, error) {
	op := "services.exercise.FetchExercisesByMuscle"
	url := fmt.Sprintf("https://api.api-ninjas.com/v1/exercises?muscle=%s", muscle)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Add("X-Api-Key", os.Getenv("API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read response: %w", op, err)
	}

	var exercises []models.ExerciseAPI
	if err := json.Unmarshal(body, &exercises); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return exercises, nil
}

func (s *ExerciseService) FetchAllExercises() ([]models.ExerciseAPI, error) {
	muscleGroups := []string{
		"abdominals", "abductors", "adductors", "biceps",
		"calves", "chest", "forearms", "glutes",
		"hamstrings", "lats", "lower_back", "middle_back",
		"neck", "quadriceps", "traps", "triceps",
	}
	var allExercises []models.ExerciseAPI
	var errs []error

	for _, muscle := range muscleGroups {
		exercises, err := s.FetchExercisesByMuscle(muscle)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed loading %s: %w", muscle, err))
			continue
		}
		allExercises = append(allExercises, exercises...)
	}
	if len(errs) > 0 {
		return allExercises, errors.Join(errs...)
	}

	return allExercises, nil
}

func (s *ExerciseService) FetchAndStoreExercises() (error) {
	op := "internal.services.fetchAndStoreExercises"
	
	exercises, err := s.FetchAllExercises()
	if err != nil {
		return fmt.Errorf("%s, failed loading exercises: %w", op, err)
	}
	if err := s.ExerciseRepo.SaveExercisesToDB(exercises); err != nil {
		return fmt.Errorf("%s, failed storing exercises: %w", op, err)
	}

	exercisesDB, err := s.ExerciseRepo.GetAllExercises()
	if err != nil{
		return fmt.Errorf("%s, failed to get all exercises: %w", op, err)
	}

	jsonString, err := json.Marshal(exercisesDB)
	if err != nil{
		return fmt.Errorf("%s, failed to marshal exercises: %w", op, err)
	}
	
	if err := s.Cache.Set(context.Background(), "exercises:all", jsonString, 0).Err(); err != nil{
		return fmt.Errorf("%s, failed to set exercises in the redis instance:  %w", op, err)
	} 
	
	for _, exercise := range exercisesDB{
		key := fmt.Sprintf("exercises:search:id:%v", exercise.ID)

		jsonString, err = json.Marshal(exercise)
		if err != nil{
			return fmt.Errorf("%s, failed to marshal exercise %s: %w", op, exercise.Name, err)
		}

		if err := s.Cache.Set(context.Background(), key, jsonString, 0).Err(); err != nil{
			return fmt.Errorf("%s, failed to set exercise %s in the redis instance:  %w", op, exercise.Name, err)
		} 
	}

	return nil
}

func (s *ExerciseService) GetAllExercises(ctx context.Context) ([]models.Exercise, error) {
    op := "internal.services.GetAllExercises"

    data, err := s.Cache.Get(ctx, "exercises:all").Result()
    if err == nil {
        var ex []models.Exercise
        if err := json.Unmarshal([]byte(data), &ex); err == nil {
            return ex, nil
        }
    } else if err != redis.Nil {
        log.Printf("warning: failed to get all exercises from redis: %v", err)
    }

	exercises, err := s.ExerciseRepo.GetAllExercises()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    payload, err := json.Marshal(exercises)
    if err != nil {
        return exercises, nil
    }
    if setErr := s.Cache.Set(ctx, "exercises:all", payload, 0).Err(); setErr != nil {
		log.Printf("warning: failed to set exercises in Redis: %v", setErr)
    }

    return exercises, nil
}

func (s *ExerciseService) GetExercisesByID(param ...interface{}) (interface{}, error){
	op := "internal.services.GetExercisesByID"

	idStr, ok := param[0].(string)
	if !ok{
        return nil, fmt.Errorf("expected id to be string, got %T", param[0])
	}

	id, err := strconv.Atoi(idStr)
	if err != nil{
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	if id == 0{
		return nil, fmt.Errorf("%s: id can not be 0", op)
	}

	key := fmt.Sprintf("exercises:search:id:%v", id)
	data, err := s.Cache.Get(context.Background(), key).Result()
	if err == nil{
		var ex models.Exercise
		if err := json.Unmarshal([]byte(data), &ex); err == nil{
			return ex, nil
		}
	}else if err != redis.Nil{
        log.Printf("warning: failed to get all exercises from redis: %v", err)
	}
	exercise, err := s.ExerciseRepo.GetExercisesByID(id)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	payload, err := json.Marshal(exercise)
    if err != nil {
        return exercise, nil
    }
	key = fmt.Sprintf("exercises:search:id:%v", exercise.ID)
	
	if err := s.Cache.Set(context.Background(), key, payload, 0).Err(); err != nil{
		log.Printf("warning: failed to set exercise %s in Redis: %v", exercise.Name, err)
	}

	return exercise, nil
}

func (s *ExerciseService) GetExercisesByName(param ...interface{}) (interface{}, error){
	op := "internal.servises.GetExercisesByName"

	name, ok := param[0].(string)
	if !ok{
        return nil, fmt.Errorf("expected name to be string, got %T", param[0])
	}

	exercise, err := s.ExerciseRepo.GetExercisesByName(name)
	if err != nil{
		return models.Exercise{}, fmt.Errorf("%s: %w", op, err)
	}

	return exercise, nil
}

func (s *ExerciseService) GetExercisesByType(param ...interface{}) (interface{}, error){
	op := "internal.servises.GetExercisesByType"

	typeEx, ok := param[0].(string)
	if !ok{
        return nil, fmt.Errorf("expected type to be string, got %T", param[0])
	}

	exercises, err := s.ExerciseRepo.GetExercisesByType(typeEx)
	if err != nil{
		return models.Exercise{}, fmt.Errorf("%s: %w", op, err)
	}

	return exercises, nil
}

func (s *ExerciseService) GetExercisesByMuscleGroup(param ...interface{}) (interface{}, error){
	op := "internal.servises.GetExercisesByMuscleGroup"

	muscleGroup, ok := param[0].(string)
	if !ok{
        return nil, fmt.Errorf("expected muscle to be string, got %T", param[0])
	}

	exercises, err := s.ExerciseRepo.GetExercisesByMuscleGroup(muscleGroup)
	if err != nil{
		return models.Exercise{}, fmt.Errorf("%s: %w", op, err)
	}

	return exercises, nil
}

func (s *ExerciseService) GetExercisesByDifficulty(param ...interface{}) (interface{}, error){
	op := "internal.servises.GetExercisesByDifficulty"

	difficulty, ok := param[0].(string)
	if !ok{
        return nil, fmt.Errorf("expected difficulty to be string, got %T", param[0])
	}

	exercises, err := s.ExerciseRepo.GetExercisesByDifficulty(difficulty)
	if err != nil{
		return models.Exercise{}, fmt.Errorf("%s: %w", op, err)
	}

	return exercises, nil
}