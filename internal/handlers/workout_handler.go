package handlers

import (
	"net/http"
	"strconv"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateWorkoutHandler(s *services.WorkoutService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var workoutCreate models.RequestCreateWorkout
		
		if err := ctx.ShouldBindJSON(&workoutCreate); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		}

		userIdRaw, exist := ctx.Get("userID")
		if !exist{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, ok := userIdRaw.(int)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not an integer"})
			return
		}

		workoutID, err := s.CreateWorkout(userID, workoutCreate)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, workoutID)
	}
}

func GetWorkoutHandler(s *services.WorkoutService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
		workoutIdStr := ctx.Query("id")
		if len(workoutIdStr) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "workout id is not set"})
		}

		workoutID, err := strconv.Atoi(workoutIdStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		}

		userID := ctx.GetInt("userID")
		
		workout, err := s.GetWorkout(workoutID, userID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, workout)
	}
}