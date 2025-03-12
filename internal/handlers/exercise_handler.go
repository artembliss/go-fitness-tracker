package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/artembliss/go-fitness-tracker/internal/models"
)


func FetchExercisesByMuscle(muscle string) ([]models.ExerciseAPI, error) {
	op := "handlers.exercises.FetchExercisesByMuscle"
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

func FetchAllExercises() ([]models.ExerciseAPI, error) {
	muscleGroups := []string{
		"abdominals", "abductors", "adductors", "biceps", 
		"calves", "chest", "forearms", "glutes", 
		"hamstrings", "lats", "lower_back", "middle_back", 
		"neck", "quadriceps", "traps", "triceps",
	}
	var allExercises []models.ExerciseAPI
	var errs []error

	for _, muscle := range muscleGroups {
		exercises, err := FetchExercisesByMuscle(muscle)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed loading %s: %w", muscle, err))
			continue
		}
		allExercises = append(allExercises, exercises...)
	}
	if len(errs) > 0 {
		return allExercises, errors.Join(errs...)
	}
	
	return allExercises, nil}

