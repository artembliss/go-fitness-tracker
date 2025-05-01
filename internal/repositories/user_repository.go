package repositories

import (
	"fmt"

	"github.com/artembliss/go-fitness-tracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (s *UserRepository) RegisterUserRepository(user models.User) (int, error){
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

func (s *UserRepository) GetUserByEmail(email string) (*models.User, error){
	const op = "repositories.user_repository.GetUserByEmail"
	
	var user models.User
	getUserQuery := `SELECT * FROM users WHERE email = $1`
	err := s.db.Get(&user, getUserQuery, email)
	if err != nil{
		return nil, fmt.Errorf("%s: failed to find user by email: %w", op, err)
	}
	return &user, nil
}

func (r *UserRepository) DeleteUser(email string, userID int) (int, error){
	const op = "repositories.DeleteUser"
	var deletedID int
	
	query := `DELETE FROM users WHERE email = $1 AND id = $2 RETURNING id`

	if err := r.db.QueryRow(query, email, userID).Scan(&deletedID); err != nil{
		return 0, fmt.Errorf("%s: failed to delete user by email: %w", op, err)
	}

	return deletedID, nil
}