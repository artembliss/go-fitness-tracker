// @title Fitness Tracker API
// @version 1.0
// @description REST API for managing workouts, programs, and users.
// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @type http
// @scheme bearer
// @in header
// @name Authorization
// @type http
// @scheme bearer

package main

import (
	"github.com/artembliss/go-fitness-tracker/internal/app"
)

func main() {
	application := &app.App{}
	application.Start()
}