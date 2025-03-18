package handlers

import (
	"net/http"
	"strconv"

	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)


func GetAllExercisesHandler(s *services.ExerciseService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		exercises, err := s.GetAllExercises()
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, &exercises)
	}
}

func GetExerciseByIdHandler(s *services.ExerciseService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var id int

		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		exercises, err := s.GetExercisesByID(id)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, exercises)
	}
}

