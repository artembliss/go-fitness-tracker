package handlers

import (
	"net/http"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(s *services.UserService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var reqUser models.RequestCreateUser
		if err := ctx.ShouldBindJSON(&reqUser); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		user, err := s.RegisterUserService(&reqUser)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"id": user.ID})
	}
}