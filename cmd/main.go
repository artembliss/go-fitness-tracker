package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/artembliss/go-fitness-tracker/internal/handlers"
	"github.com/artembliss/go-fitness-tracker/logger/sl"
	"github.com/artembliss/go-fitness-tracker/storage/postgre"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	log := sl.SetUpLogger(env)
	log.Info("Starting server", slog.String("env", env))

	storage, err := postgre.New()
	if err != nil{
		log.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	log.Info("Storage initialized")

	if handlers.CheckExercisesExist(&storage){
		log.Info("Exercises exist")
	}else{
		if err := fetchAndStoreExercises(&storage); err != nil{
			log.Error("Failed to fetch exercises", sl.Err(err))
		}
		log.Info("Exercises fetched successfully")
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env: %s", err)
	}
}

func fetchAndStoreExercises(storage *postgre.Storage) (error) {
	op := "main.fetchAndStoreExercises"
	exercises, err := handlers.FetchAllExercises()
	if err != nil {
		return fmt.Errorf("%s, failed loading exercises: %w", op, err)
	}
	return handlers.SaveExercisesToDB(storage, exercises)	
}