package handlers

import (
	"net/http"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

// RegisterUserHandler godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.RequestCreateUser true "User registration information"
// @Success 200 {object} map[string]int "Registered user ID"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/register [post]
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

// LoginUserHandler godoc
// @Summary Authenticate user and get token
// @Description User login to obtain JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.RequestLoginUser true "User login credentials"
// @Success 200 {object} map[string]string "JWT Token"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/login [post]
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

// GetUserHandler godoc
// @Summary Get user by email
// @Description Retrieve user information using email address
// @Security BearerAuth
// @Tags Users
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user [get]
func GetUserHandler(s *services.UserService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		email := ctx.Query("email")

		if len(email) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user email is not set"})
			return
		}

		user, err := s.GetUserByEmail(email)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

// DeleteUserHandler godoc
// @Summary Delete user by email
// @Description Delete a user account using their email address
// @Security BearerAuth
// @Tags Users
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Success 200 {integer} int "Deleted user ID"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user [delete]
func DeleteUserHandler(s *services.UserService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		email := ctx.Query("email")

		if len(email) == 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user email is not set"})
			return
		}

		userID := ctx.GetInt("userID") 

		deletedID, err := s.DeleteUser(email, userID)
		if err != nil{
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, deletedID)
	}
}