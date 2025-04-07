package handlers

import (
	"net/http"
	"strconv"

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

		createdID, err := s.CreateProgram(program)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})		
			return
		}
		
		ctx.JSON(http.StatusOK, createdID)
	}
}

func GetProgramHandler(s *services.ProgramService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		userID := ctx.GetInt("userID")

		programIdStr := ctx.Query("id")
		if len(programIdStr) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "program id not set"})
		}

		programID, err := strconv.Atoi(programIdStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		}

		programs, err := s.GetProgram(programID, userID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, programs)
	}
}

func UpdateProgramHandler(s *services.ProgramService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var programUpdate models.RequestCreateProgram
		
		if err := ctx.ShouldBindJSON(&programUpdate); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})		
			return	
		}

		idStr := ctx.Query("id")
		programID, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		nameToID, err := s.GetNameToID(programUpdate.Exercises)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find exercises in storage"})
			return
		}

		exercisesToSave, notFound := s.MapToDBExercises(programUpdate.Exercises, nameToID)
		if len(notFound) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "some exercises not found"})		
			return
		}

		program := models.Program{
			UserID: ctx.GetInt("userID"),
			Name: programUpdate.Name,
			Exercises: exercisesToSave,	
		}

		ID, err := s.UpdateProgram(program, programID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})		
			return
		}
		
		ctx.JSON(http.StatusOK, ID)
	}
}

//Delete row from exercises_program
func DeleteProgramHandler(s *services.ProgramService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
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

		idStr := ctx.Query("id")
		programID, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		deletedID, err := s.DeleteProgram(programID, userID)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, deletedID)
	}
}