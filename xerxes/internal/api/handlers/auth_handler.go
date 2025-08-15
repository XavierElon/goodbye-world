package handlers

import (
	"encoding/json"
	"net/http"

	"go-api/internal/domain"
	"go-api/internal/services"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// SendVerificationCode godoc
// @Summary      Send verification code
// @Description  Sends an SMS verification code to the user's phone number
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body domain.VerificationRequest true "Phone number"
// @Success      200 {object} domain.VerificationResponse
// @Failure      400 {string} string "Invalid request"
// @Failure      500 {string} string "Internal error"
// @Router       /auth/send-code [post]
func (h *AuthHandler) SendVerificationCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.PhoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}
	resp, err := h.Service.SendVerificationCode(req.PhoneNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// VerifyCodeAndLogin godoc
// @Summary      Verify code and login/register
// @Description  Verifies the code and logs in or registers the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body domain.UserLogin true "Phone number and code"
// @Success      200 {object} domain.AuthResponse
// @Failure      400 {string} string "Invalid request"
// @Failure      401 {string} string "Unauthorized"
// @Router       /auth/verify [post]
func (h *AuthHandler) VerifyCodeAndLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.PhoneNumber == "" || req.Code == "" {
		http.Error(w, "Phone number and code are required", http.StatusBadRequest)
		return
	}
	resp, err := h.Service.VerifyCodeAndLogin(req.PhoneNumber, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}