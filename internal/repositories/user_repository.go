package repositories

import (
	"fmt"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) RegisterUserRepository(user models.User) (int, error){
	const op = "repositories.RegisterUserRepository" 
	query := `INSERT INTO users (name, email, password_hash, age, gender, height, weight, created_at) 
			VALUES (:name, :email, :password_hash, :age, :gender, :height, :weight, NOW()) RETURNING id`
	rows, err := s.db.NamedQuery(query, user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
		}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&user.ID); err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}
	return user.ID, nil
}