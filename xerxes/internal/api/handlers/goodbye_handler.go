package handlers

import (
	"go-api/internal/services"
	"net/http"
)

type GoodbyeHandler struct {
    Service *services.GoodbyeService
}

func NewGoodbyeHandler(service *services.GoodbyeService) *GoodbyeHandler {
    return &GoodbyeHandler{Service: service}
}

// @Summary Goodbye world endpoint
// @Description Returns a goodbye message
// @Tags goodbye
// @Accept json
// @Produce plain
// @Success 200 {string} string "Goodbye, World!"
// @Router /goodbyeworld [get]
func (h *GoodbyeHandler) GoodbyeWorld(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(h.Service.GoodbyeWorld()))
}