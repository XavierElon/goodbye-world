package routes

import (
	"go-api/internal/api/handlers"
	"go-api/internal/services"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func RegisterRoutes() {
	healthService := services.NewHealthService()
	healthHandler := handlers.NewHealthHandler(healthService)

	goodbyeService := services.NewGoodbyeService()
    goodbyeHandler := handlers.NewGoodbyeHandler(goodbyeService)

    authService := services.NewAuthService()
    authHandler := handlers.NewAuthHandler(authService)

    http.HandleFunc("/health", healthHandler.Health)
    http.HandleFunc("/goodbyeworld", goodbyeHandler.GoodbyeWorld)
    http.HandleFunc("/auth/send-code", authHandler.SendVerificationCode)
    http.HandleFunc("/auth/verify", authHandler.VerifyCodeAndLogin)
    
    // Serve swagger.json directly
    http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        http.ServeFile(w, r, "docs/swagger.json")
    })
    
    // Swagger documentation
    http.HandleFunc("/swagger/", httpSwagger.Handler(
        httpSwagger.URL("http://localhost:10471/swagger/doc.json"), // The URL pointing to API definition
    ))
}