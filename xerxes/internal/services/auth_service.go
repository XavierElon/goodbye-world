package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"go-api/internal/config"
	"go-api/internal/domain"
	"go-api/internal/repository"
)

type AuthService struct {
	repo *repository.RedisRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repository.NewRedisRepository(),
	}
}

// SendVerificationCode sends an SMS verification code to the user's phone
func (s *AuthService) SendVerificationCode(phoneNumber string) (*domain.VerificationResponse, error) {
	// Generate a 6-digit verification code
	code := generateVerificationCode()
	
	// Store the code in Redis (expires in 10 minutes)
	ctx := context.Background()
	err := s.repo.StoreVerificationCode(ctx, phoneNumber, code)
	if err != nil {
		return nil, fmt.Errorf("failed to store verification code: %w", err)
	}
	
	// Send SMS via AWS SNS
	message := fmt.Sprintf("Your verification code is: %s. Valid for 10 minutes.", code)
	err = config.SendSMS(phoneNumber, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send SMS: %w", err)
	}
	
	log.Printf("Verification code sent to %s: %s", phoneNumber, code)
	
	return &domain.VerificationResponse{
		Message: "Verification code sent successfully",
		Success: true,
	}, nil
}

// VerifyCodeAndLogin verifies the code and logs the user in
func (s *AuthService) VerifyCodeAndLogin(phoneNumber, code string) (*domain.AuthResponse, error) {
	ctx := context.Background()
	
	// Get the stored verification code
	storedCode, err := s.repo.GetVerificationCode(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired verification code: %w", err)
	}
	
	// Check if codes match
	if storedCode != code {
		return nil, fmt.Errorf("invalid verification code")
	}
	
	// Get or create user
	user, err := s.repo.GetUser(ctx, phoneNumber)
	if err != nil {
		// User doesn't exist, create new one
		user = &domain.User{
			ID:          generateUserID(),
			PhoneNumber: phoneNumber,
			CreatedAt:   time.Now(),
			LastLogin:   time.Now(),
			IsVerified:  true,
		}
		
		// Store the new user
		err = s.repo.StoreUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update last login
		user.LastLogin = time.Now()
		user.IsVerified = true
		
		// Update user in Redis
		err = s.repo.StoreUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}
	
	// Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	// Clear the verification code after successful login
	s.repo.ClearVerificationCode(ctx, phoneNumber)
	
	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// generateVerificationCode creates a random 6-digit code
func generateVerificationCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	
	// Convert to 6-digit number
	code := int(bytes[0])<<16 | int(bytes[1])<<8 | int(bytes[2])
	code = code % 1000000 // Ensure it's 6 digits
	
	// Format as string with leading zeros if needed
	return fmt.Sprintf("%06d", code)
}

// generateUserID creates a unique user ID
func generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// generateJWTToken creates a JWT token for the user
func generateJWTToken(user *domain.User) (*domain.AuthToken, error) {
	// For now, we'll create a simple token structure
	// In the next step, we'll implement proper JWT
	
	expiresAt := time.Now().Add(24 * time.Hour)
	
	token := &domain.AuthToken{
		AccessToken: fmt.Sprintf("token_%s_%d", user.ID, expiresAt.Unix()),
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 24 hours in seconds
		ExpiresAt:   expiresAt,
	}
	
	return token, nil
}