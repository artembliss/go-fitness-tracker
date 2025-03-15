package services

import (
	"fmt"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/artembliss/go-fitness-tracker/pkg/auth"
)
type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (s *UserService) RegisterUserService(reqUser *models.RequestCreateUser) (models.User, error){
	const op = "services.RegisterUserService"
	
	hashedPassword, err := auth.HashPassword(reqUser.Password)
	if err != nil{
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user := models.User{
		Name: reqUser.Name,
		Email: reqUser.Email,
		PasswordHash: hashedPassword,
		Age: reqUser.Age,
		Gender: reqUser.Gender,
		Height: reqUser.Height,
		Weight: reqUser.Weight,
	}

	userID, err := s.UserRepo.RegisterUserRepository(user)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	user.ID = userID

	return user, nil
}