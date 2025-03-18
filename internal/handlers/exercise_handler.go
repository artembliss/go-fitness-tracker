package handlers

import (
	"net/http"

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


