package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/artembliss/go-fitness-tracker/internal/handlers"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/artembliss/go-fitness-tracker/logger/sl"
	"github.com/artembliss/go-fitness-tracker/storage/postgre"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	logger := sl.SetUpLogger(env)
	logger.Info("Starting server", slog.String("env", env))

	storage, err := postgre.New()
	if err != nil{
		logger.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	logger.Info("Storage initialized")

	if repositories.CheckExercisesExist(&storage){
		logger.Info("Exercises exist")
	}else{
		if err := fetchAndStoreExercises(&storage); err != nil{
			logger.Error("Failed to fetch exercises", sl.Err(err))
		}
		logger.Info("Exercises fetched successfully")
	}
	userStorage := repositories.NewUserRepository(storage.GetDB())
	userService := services.NewUserService(userStorage)
	
	authService := services.NewAuthService(userStorage)

	router := gin.Default()
  
	router.POST("/register", handlers.RegisterUserHandler(userService))
	router.POST("/login", handlers.LoginUserHandler(authService))

  
	if err := router.Run(":8080"); err != nil {
	  logger.Error("Failed to start server:", sl.Err(err))
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
	return repositories.SaveExercisesToDB(storage, exercises)	
}