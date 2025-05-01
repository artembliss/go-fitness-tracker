package handlers

import (
	"net/http"

	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

// GetAllExercisesHandler godoc
// @Summary Get all exercises
// @Description Retrieve a list of all available exercises
// @Tags Exercises
// @Accept json
// @Produce json
// @Success 200 {array} models.Exercise
// @Failure 404 {object} map[string]string
// @Router /exercises [get]
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

// GetExerciseByParamHandler godoc
// @Summary Search exercises by parameter
// @Description Find exercises by one of the following parameters: id, name, type, muscle group, or difficulty
// @Tags Exercises
// @Accept json
// @Produce json
// @Param id query string false "Exercise ID"
// @Param name query string false "Exercise name"
// @Param type query string false "Exercise type"
// @Param muscle query string false "Target muscle group"
// @Param difficulty query string false "Exercise difficulty level"
// @Success 200 {array} models.Exercise
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /exercises/search [get]
func GetExerciseByParamHandler(s *services.ExerciseService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		requestParams := make(map[string]string)

		servicesByParam := map[string]services.ServiceFunc{
			"id": s.GetExercisesByID, "name": s.GetExercisesByName, "type": s.GetExercisesByType,
		    "muscle": s.GetExercisesByMuscleGroup, "difficulty": s.GetExercisesByDifficulty}

		paramNames := []string{"id", "name", "type", "muscle", "difficulty"}
		
		for _, paramName := range(paramNames){
			param := ctx.Query(paramName)
			if param != ""{
				requestParams[paramName] = param
			}
		}

		if len(requestParams) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "parameters not set"})
			return
		}
		if len(requestParams) > 1{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "need to be only 1 parameter"})
			return
		}

		for name, param := range requestParams {
			serviceFunc, ok := servicesByParam[name]
			if !ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported parameter " + name})
				return
			}
			result, err := serviceFunc(param)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, result)
			return
		}
	}
}

