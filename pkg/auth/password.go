package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pswd string) (string, error) {
	const op = "services.password.HashPassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil{
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(pswd, hash string) (bool){
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pswd))
	return err == nil
}