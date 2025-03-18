package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/artembliss/go-fitness-tracker/internal/handlers"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/artembliss/go-fitness-tracker/pkg/logger/sl"
	"github.com/artembliss/go-fitness-tracker/pkg/storage/postgre"
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

	UserRepository := repositories.NewUserRepository(storage.GetDB())
	ExerciseRepository := repositories.NewExerciseRepository(storage.GetDB())

	userService := services.NewUserService(UserRepository)
	authService := services.NewAuthService(UserRepository)
	exerciseService := services.NewExerciseService(ExerciseRepository)

	logger.Info("Storage initialized")

	if ExerciseRepository.CheckExercisesExist(){
		logger.Info("Exercises exist")
	}else{
		if err := exerciseService.FetchAndStoreExercises() ; err != nil{
			logger.Error("Failed to fetch exercises", sl.Err(err))
		}
		logger.Info("Exercises fetched successfully")
	}


	router := gin.Default()
  
	router.POST("/register", handlers.RegisterUserHandler(userService))
	router.POST("/login", handlers.LoginUserHandler(authService))

	router.GET("/exercises", handlers.GetAllExercisesHandler(exerciseService))


	if err := router.Run(":8080"); err != nil {
	  logger.Error("Failed to start server:", sl.Err(err))
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env: %s", err)
	}
}

