package handlers

import (
	"net/http"
	"os"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

var jwtSecretKey = []byte(os.Getenv("JWT_KEY"))


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

func LoginUserHandler(s *services.AuthService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var userLogin models.RequestLoginUser

		if err := ctx.ShouldBindJSON(&userLogin); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		token, err := s.AuthenticateUserService(userLogin.Email, userLogin.Password)
		if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
            return
        }

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}