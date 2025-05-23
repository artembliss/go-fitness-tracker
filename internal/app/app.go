package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"

	_ "github.com/artembliss/go-fitness-tracker/docs"
	"github.com/artembliss/go-fitness-tracker/internal/handlers"
	"github.com/artembliss/go-fitness-tracker/internal/middleware"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/artembliss/go-fitness-tracker/pkg/logger/sl"
	"github.com/artembliss/go-fitness-tracker/pkg/storage/postgre"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type App struct {
	router *gin.Engine
	logger *slog.Logger
	db     *postgre.Storage
	cache  *redis.Client
}

func (a *App) InitConfig(){
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error uploading .env: %s", err)
	}
}

func (a *App) InitLogger(){
	env := os.Getenv("ENV")
	a.logger = sl.SetUpLogger(env)
	a.logger.Info("Starting server", slog.String("env", env))
}

func (a *App) InitDB(){
	storage, err := postgre.New()
	if err != nil {
		a.logger.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}
	a.logger.Info("Storage initialized")
	
	a.db = storage
}

func (a *App) InitRedis(){
	addr := os.Getenv("REDIS_ADDR")
    dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil{
		a.logger.Error("failed to fetch REDIS_DB .env var", sl.Err(err))
		os.Exit(1)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
        Password: "", 
        DB:       dbNum,  
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil{
		a.logger.Error("failed to init cashe", sl.Err(err))
		os.Exit(1)
	}

	a.cache = rdb
}


func (a *App) InitRouters(storage *postgre.Storage, cache *redis.Client) {
	db := storage.GetDB()

	userRepo := repositories.NewUserRepository(db)
	exerciseRepo := repositories.NewExerciseRepository(db)
	programRepo := repositories.NewProgramRepository(db)
	workoutRepo := repositories.NewWorkoutRepository(db)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	exerciseService := services.NewExerciseService(exerciseRepo, cache)
	programService := services.NewProgramService(programRepo)
	workoutService := services.NewWorkoutService(workoutRepo)

	authMiddleware := middleware.JWTMiddleware(userService)

	if exerciseRepo.CheckExercisesExist() {
		a.logger.Info("Exercises exist")
	} else {
		if err := exerciseService.FetchAndStoreExercises(); err != nil {
			a.logger.Error("Failed to fetch exercises", sl.Err(err))
		} else {
			a.logger.Info("Exercises fetched successfully")
		}
	}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/user/register", handlers.RegisterUserHandler(userService))
	router.POST("/user/login", handlers.LoginUserHandler(authService))

	router.GET("/exercises", handlers.GetAllExercisesHandler(exerciseService))
	router.GET("/exercises/search", handlers.GetExerciseByParamHandler(exerciseService))

	protected := router.Group("/", authMiddleware)
	{
		protected.GET("/user", handlers.GetUserHandler(userService))
		protected.DELETE("/user", handlers.DeleteUserHandler(userService))

		protected.POST("/programs", handlers.CreateProgramHandler(programService))
		protected.GET("/programs", handlers.GetProgramHandler(programService))
		protected.DELETE("/programs", handlers.DeleteProgramHandler(programService))
		protected.PATCH("/programs", handlers.UpdateProgramHandler(programService))

		protected.POST("/workouts", handlers.CreateWorkoutHandler(workoutService))
		protected.GET("/workouts", handlers.GetWorkoutHandler(workoutService))
		protected.DELETE("/workouts", handlers.DeleteWorkoutHandler(workoutService))
		protected.PATCH("/workouts", handlers.UpdateWorkoutHandler(workoutService))
	}

	a.router = router
}

func (a *App) Start() {
	a.InitConfig()
	a.InitLogger()
	a.InitDB()
	a.InitRedis()
	a.InitRouters(a.db, a.cache)

	if err := a.router.Run(":8080"); err != nil {
		a.logger.Error("Failed to start server:", sl.Err(err))
	}
}