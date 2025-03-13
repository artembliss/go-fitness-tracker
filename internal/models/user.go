package models

import "time"

type User struct {
	ID           int    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Age          int       `json:"age" db:"age"`
	Gender       string    `json:"gender" db:"gender"`
	Height       int       `json:"height" db:"height"`
	Weight       float64   `json:"weight" db:"weight"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
}

type RequestCreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Height   int    `json:"height"`
	Weight   float64 `json:"weight"`
}

type RequestUpdateUser struct {
	ID       int  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Age      int     `json:"age"`
	Gender   string  `json:"gender"`
	Height   int     `json:"height"`
	Weight   float64 `json:"weight"`
}

type RequestLoginUser struct{
	Email    string  `json:"email"`
	Password string  `json:"password"`
}

