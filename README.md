# go‑fitness‑tracker <!-- omit in toc -->
Fitness Tracker is a Go (Gin) REST API for tracking programs, workouts and exercises. 
It imports data into PostgreSQL, leverages Redis for caching, supports JWT-based auth, 
provides Swagger documentation, the project can be run either using Docker or directly on your local machine.

- [🚀 About the project](#-about-the-project)
- [📌 Core features](#-core-features)
- [🛠 Tech stack](#-tech-stack)
- [🗂 Project structure](#project-structure)
- [🗺 Quick start](#quick-start)
  
---

## 🚀 About the project
Fitness Tracker lets users:

1. **Register / log in** (JWT; bcrypt password hashing).
2. **Create workout programs** composed of reference exercises.
3. **Log workouts** with per‑set reps & weights.
4. Retrieve the **exercise catalogue** — cached in Redis.

> **Status:** MVP v1.1 (Redis cache for `GET /exercises` & `GET /exercises/:id` is live).

---

## 📌 Core features
- **JWT Auth** — protected routes via middleware.  
- **Exercise catalogue** (≈1400 movements) with external API import.
- **CRUD operations** for workouts and programs only for ouners
- **Search exercises** by name, muscle group, difficulty and other parameters
- **Redis** look‑aside cache on read‑heavy catalogue queries.  
- **Swagger UI** (`/swagger/index.html`).  

---

## 🛠 Tech stack
| Layer     | Tech                        |
|-----------|-----------------------------|
| Backend   | Go 1.22, Gin                |
| Database  | PostgreSQL 16, sqlx         |
| Cache     | Redis 7 (go‑redis v9)       |
| Auth      | JWT, bcrypt                 |
| Docs      | Swag ▶︎ Swagger UI          |
| DevOps    | Docker Compose v3.8         |

---

## 🗂 Project structure
```text
fitness-tracker
├── cmd/                     # CLI / entry‑points
│   └── main.go              # starts the HTTP server
├── docs/                    # Swagger / OpenAPI files
├── internal/                # Application logic
│   ├── app/                 # InitConfig, Start(), DI bootstrap
│   ├── handlers/            # Gin HTTP handlers (controllers)
│   ├── middleware/          # Auth, logging, error recovery
│   ├── models/              # Domain data models
│   ├── repositories/        # SQLx queries & persistence
│   └── services/            # Core business rules
├── pkg/                     # Reusable packages
│   ├── auth/                # JWT helpers & password hashing
│   ├── logger/              # slog wrappers
│   ├── migrations/          # SQL migration files
│   └── storage/postgre/     # Postgres connection & helpers
├── .env                     # Runtime secrets (ignored in VCS)
├── .env.example             # Sample env config
├── .gitignore
├── docker-compose.yaml      # Local stack: Postgres, Redis, migrate
├── go.mod
├── go.sum
└── README.md                # You are here
```

---
## 🗺 Quick start

---

### 1. Prerequisites
* Go ≥ 1.22 *(only if you want local `go run`)*  
* Docker Engine / Docker Desktop  

--- 
### 2. Clone the repository
```bash
git clone https://github.com/artembliss/fitness-tracker.git
cd fitness-tracker
```
### 2. Setup & Installation
### Prerequisites
- Install [Go](https://go.dev/doc/install)
- Install [Docker](https://docs.docker.com/get-docker/)

---

### 3. Configure environment
```bash
cp .env.example .env
```

### 4. Run the application
```sh
git clone https://github.com/artembliss/fitness-tracker.git
cd fitness-tracker
```

### Run with Docker
```sh
docker-compose up --build
```
---

### Run Locally
1. Copy the `.env.example` file and configure environment variables:
```sh
cp .env.example .env
```

2. Make sure you have PostgreSQL running (you can do this via Docker)
 ```bash
  docker run -d \
  --name ft-postgres \
  -e POSTGRES_USER=${DB_USER} \
  -e POSTGRES_PASSWORD=${DB_PASSWORD} \
  -e POSTGRES_DB=${DB_NAME} \
  -p 5432:5432 \
  postgres:16.8-alpine3.20
```
3. Install dependencies
  ```sh
go mod tidy
``` 
4. Run the server:
```sh
go run cmd/main.go
```

---
### API Documentation - Swagger UI
Access interactive API documentation at:
```sh
http://localhost:8080/swagger/index.html
```
![image](https://github.com/user-attachments/assets/7de08655-3bc3-46fd-bf59-452cf6ba391c)

