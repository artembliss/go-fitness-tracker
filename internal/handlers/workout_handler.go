package handlers

import (
	"net/http"
	"strconv"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateWorkoutHandler godoc
// @Summary Create a new workout
// @Description Create a new workout for the authenticated user
// @Security BearerAuth
// @Tags Workouts
// @Accept json
// @Produce json
// @Param workout body models.RequestCreateWorkout true "Workout information"
// @Success 200 {integer} int "Created workout ID"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /workouts [post]
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

// GetWorkoutHandler godoc
// @Summary Get a workout by ID
// @Description Retrieve a specific workout for the authenticated user
// @Security BearerAuth
// @Tags Workouts
// @Accept json
// @Produce json
// @Param id query int true "Workout ID"
// @Success 200 {object} models.Workout
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /workouts [get]
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

// DeleteWorkoutHandler godoc
// @Summary Delete a workout by ID
// @Description Delete a specific workout for the authenticated user
// @Security BearerAuth
// @Tags Workouts
// @Accept json
// @Produce json
// @Param id query int true "Workout ID"
// @Success 200 {integer} int "Deleted workout ID"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /workouts [delete]
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

// UpdateWorkoutHandler godoc
// @Summary Update an existing workout
// @Description Update a specific workout's details for the authenticated user
// @Security BearerAuth
// @Tags Workouts
// @Accept json
// @Produce json
// @Param id query int true "Workout ID"
// @Param workout body models.RequestCreateWorkout true "Updated workout information"
// @Success 200 {integer} int "Updated workout ID"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /workouts [patch]
func UpdateWorkoutHandler(s *services.WorkoutService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var workoutUpdate models.RequestCreateWorkout
		
		if err := ctx.ShouldBindJSON(&workoutUpdate); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		
		idStr := ctx.Query("id")
		if len(idStr) == 0{
			ctx.JSON(http.StatusBadRequest,  gin.H{"error": "workout id is not set"})
			return
		}
		
		id, err := strconv.Atoi(idStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest,  gin.H{"error": "invalid workout id"})
			return
		}

		userID := ctx.GetInt("userID") 

		updatedID, err := s.UpdateWorkout(id, userID, workoutUpdate)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, updatedID)
		
	}
}