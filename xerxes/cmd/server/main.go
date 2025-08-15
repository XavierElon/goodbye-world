package main

import (
	"go-api/internal/api/routes"
	"go-api/internal/config"
	"log"
	"net/http"
)

// @title Go API
// @version 1.0
// @description A simple Go API with health and goodbye endpoints
// @host localhost:10471
// @BasePath /
func main() {
	// Initialize Redis
	config.InitRedis()

	// Initialize AWS
	config.InitAWS()

    routes.RegisterRoutes()

    log.Println("Server running on http://localhost:10471")
	log.Println("Health check: http://localhost:10471/health")
	log.Println("Goodbye world: http://localhost:10471/goodbyeworld")
	log.Println("Auth send code: http://localhost:10471/auth/send-code")
	log.Println("Auth verify: http://localhost:10471/auth/verify")
	log.Println("Swagger docs: http://localhost:10471/swagger/index.html")
	log.Println("Press Ctrl+C to stop")
	
    log.Fatal(http.ListenAndServe(":10471", nil))
}