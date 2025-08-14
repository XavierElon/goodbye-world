package handlers

import (
	"go-api/internal/services"
	"net/http"
)

type HealthHandler struct {
    Service *services.HealthService
}

func NewHealthHandler(service *services.HealthService) *HealthHandler {
    return &HealthHandler{Service: service}
}

// @Summary Health check endpoint
// @Description Returns the health status of the API
// @Tags health
// @Accept json
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(h.Service.Health()))
}