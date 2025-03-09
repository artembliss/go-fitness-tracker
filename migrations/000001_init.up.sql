	CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password TEXT NOT NULL,
	age INT,
	gender VARCHAR(20),
	height INT,
	weight FLOAT,
	created_at TIMESTAMP DEFAULT now() NOT NULL
    );

    CREATE TABLE programs(
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,  
	name VARCHAR(255) NOT NULL,
	exercises JSONB NOT NULL,
	created_at TIMESTAMP DEFAULT now() NOT NULL
	);

    CREATE TABLE workouts(
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	program_id INT REFERENCES programs(id) ON DELETE SET NULL,
	date DATE NOT NULL,
	exercises JSONB NOT NULL,
	duration INT, 
	calories FLOAT,
	created_at TIMESTAMP DEFAULT now() NOT NULL
	);

    CREATE TABLE exercises(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	type VARCHAR(255),
	muscle_group VARCHAR(50) NOT NULL,
	equipment VARCHAR(255),
	difficulty VARCHAR(255),
	instruction TEXT);