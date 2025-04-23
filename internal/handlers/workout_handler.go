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
			return
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
			return
		}

		ctx.JSON(http.StatusOK, workoutID)
	}
}

func GetWorkoutHandler(s *services.WorkoutService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
		workoutIdStr := ctx.Query("id")
		if len(workoutIdStr) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "workout id is not set"})
			return
		}

		workoutID, err := strconv.Atoi(workoutIdStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
			return
		}

		userID := ctx.GetInt("userID")
		
		workout, err := s.GetWorkout(workoutID, userID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, workout)
	}
}

func DeleteWorkoutHandler(s *services.WorkoutService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		idStr := ctx.Query("id")
		if len(idStr) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "workout id is not set"})
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
			return
		}

		userID := ctx.GetInt("userID")

		deletedWorkoutId, err := s.DeleteWorkout(id, userID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, deletedWorkoutId)
	}
}