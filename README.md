# go-fitness-tracker
Fitness tracker with training programs, exercise logging and user progress.

# Fitness Tracker

## 🚀 About the Project
Fitness Tracker is a REST API built in **Go** for tracking fitness workouts. Users can register, create workout programs, log exercises, and track progress over time.

## 📌 Features
- User authentication with JWT
- Workout program management
- Exercise tracking within workouts
- PostgreSQL as the database
- Uses **Gin** for routing and **sqlx** for database queries
- Dockerized environment for easy deployment
- Future plans: Redis caching, test coverage, analytics

## 🛠 Tech Stack
- **Go** (Gin, sqlx)
- **PostgreSQL**
- **Docker**
- **JWT** for authentication

## 🏗 Setup & Installation
### Prerequisites
- Install [Go](https://go.dev/doc/install)
- Install [Docker](https://docs.docker.com/get-docker/)

### Clone the Repository
```sh
git clone https://github.com/artembliss/fitness-tracker.git
cd fitness-tracker
```

### Run with Docker
```sh
docker-compose up --build
```

### Run Locally
1. Copy the `.env.example` file and configure environment variables:
```sh
cp .env.example .env
```
2. Start PostgreSQL (or use Docker for it).
3. Run the server:
```sh
go run main.go
```

## 📖 API Endpoints
### **Authentication**
| Method | Endpoint       | Description  |
|--------|--------------|--------------|
| POST   | `/register`  | Register a new user |
| POST   | `/login`     | Login and get JWT  |

### **Workouts**
| Method | Endpoint             | Description |
|--------|----------------------|-------------|
| GET    | `/workouts`          | Get user workouts |
| POST   | `/workouts`          | Create a new workout |
| GET    | `/workouts/:id`      | Get workout details |

## ✅ To-Do List
- [ ] Implement Redis caching
- [ ] Add unit and integration tests
- [ ] Implement analytics & user progress tracking


