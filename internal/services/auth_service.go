package services

import (
	"fmt"

	"github.com/artembliss/go-fitness-tracker/internal/repositories"
	"github.com/artembliss/go-fitness-tracker/pkg/auth"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) AuthenticateUserService(email, password string) (string, error){
	const op = "services.auth_service.AuthenticateUserService" 

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil{
		return "", fmt.Errorf("%s: user not found: %w", op, err)
	}

	if !auth.CheckPassword(password, user.PasswordHash){
		return "", fmt.Errorf("%s: Ivalid email or password: %w", op, err)
	}

	token, err := auth.GenerateJWT(user.Email) 
	if err != nil{
		return "", fmt.Errorf("%s: failed to generate token: %w", op, err)
	}

	return token, nil
}