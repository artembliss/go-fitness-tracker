package postgre

import (
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) GetDB() *sqlx.DB {
	return s.db
}


func New() (Storage, error){
	const op = "storage.postgre.New"	

	dbUser, userExist := os.LookupEnv("DB_USER")
	dbName, nameExist := os.LookupEnv("DB_NAME")
	dbPswd, pswdExist := os.LookupEnv("DB_PASSWORD")
	host, hostExist := os.LookupEnv("DB_HOST")
	port, portExist := os.LookupEnv("DB_PORT")
	if !userExist || !nameExist || !pswdExist || !hostExist || !portExist{
		return Storage{}, fmt.Errorf("%s: some env variables not set", op)
	} 
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	host, port, dbUser, dbPswd, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil{
		return Storage{}, fmt.Errorf("%s: failed to connect to storage: %w", op, err)
	}

	createTableUsersQuery := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		age INT,
		gender VARCHAR(20),
		height INT,
		weight FLOAT,
		created_at TIMESTAMP DEFAULT now() NOT NULL
	)`
	if _, err := db.Exec(createTableUsersQuery); err != nil{
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}
	
	createTableProgramsQuery := `
	CREATE TABLE IF NOT EXISTS programs(
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,  
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT now() NOT NULL
	)`
	if _, err := db.Exec(createTableProgramsQuery); err != nil{
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	createTableWorkoutsQuery := `
	CREATE TABLE IF NOT EXISTS workouts(
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	program_id INT REFERENCES programs(id) ON DELETE SET NULL,
	date DATE NOT NULL,
	exercises JSONB NOT NULL,
	duration INT, 
	calories FLOAT,
	created_at TIMESTAMP DEFAULT now() NOT NULL
	)`
	if _, err := db.Exec(createTableWorkoutsQuery); err != nil{
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	createTableExercisesQuery := `
	CREATE TABLE IF NOT EXISTS exercises(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	type VARCHAR(255),
	muscle_group VARCHAR(50),
	equipment VARCHAR(255),
	difficulty VARCHAR(255),
	instruction TEXT)`
	if _, err := db.Exec(createTableExercisesQuery); err != nil{
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	createTableExercisesProgramQuery := `
	CREATE TABLE IF NOT EXISTS exercises_program(
	id SERIAL PRIMARY KEY,
	program_id INT REFERENCES programs(id) ON DELETE CASCADE,
	exercise_id INT REFERENCES exercises(id) ON DELETE CASCADE,
	sets INTEGER NOT NULL,
	reps INTEGER NOT NULL,
	weight DECIMAL(6,3))`
	if _, err := db.Exec(createTableExercisesProgramQuery); err != nil{
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	createTableExercisesEntryQuery := `
	CREATE TABLE IF NOT EXISTS exercises_entry(
	id SERIAL PRIMARY KEY,
	workout_id INT REFERENCES workouts(id) ON DELETE CASCADE,
	exercise_id INT REFERENCES exercises(id) ON DELETE CASCADE,
	sets INT[] NOT NULL,
	reps INT[] NOT NULL,
	weight DECIMAL(6,3)[])`
	if _, err := db.Exec(createTableExercisesEntryQuery); err != nil{
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}
	return Storage{db: db}, nil
}