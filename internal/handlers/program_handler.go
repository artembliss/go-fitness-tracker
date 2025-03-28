package handlers

import (
	"net/http"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateProgramHandler(s *services.ProgramService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var programCreate models.RequestCreateProgram
		
		if err := ctx.ShouldBindJSON(&programCreate); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})		
			return	
		}

		nameToID, err := s.GetNameToID(programCreate.Exercises)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find exercises in storage"})
			return
		}

		exercisesToSave, notFound := s.MapToDBExercises(programCreate.Exercises, nameToID)
		if len(notFound) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "some exercises not found"})		
			return
		}

		program := models.Program{
			UserID: ctx.GetInt("userID"),
			Name: programCreate.Name,
			Exercises: exercisesToSave,	
		}

		createdID, err := s.SaveProgram(program)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})		
			return
		}
		
		ctx.JSON(http.StatusOK, createdID)
	}
}

func GetProgramsHandler(s *services.ProgramService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		userID := ctx.GetInt("userID")

		programs, err := s.GetPrograms(userID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, programs)
	}
}