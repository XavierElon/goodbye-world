package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-api/internal/config"
	"go-api/internal/domain"
)

type RedisRepository struct{}

func NewRedisRepository() *RedisRepository {
	return &RedisRepository{}
}

func (r *RedisRepository) StoreVerificationCode(ctx context.Context, phoneNumber string, code string) error {
	key := fmt.Sprintf("verification:%s", phoneNumber)

	// Store code with 10 minute TTL
	err := config.RedisClient.Set(ctx, key, code, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	return nil
}

func (r *RedisRepository) GetVerificationCode(ctx context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf("verification:%s", phoneNumber)

	code, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("verification code not found or expired: %w", err)
	}
	
	return code, nil
}

func (r *RedisRepository) ClearVerificationCode(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("verification:%s", phoneNumber)
	
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to clear verification code: %w", err)
	}
	
	return nil
}

func (r *RedisRepository) StoreUser(ctx context.Context, user *domain.User) error {
	// Store user profile
	userKey := fmt.Sprintf("user:profile:%s", user.PhoneNumber)
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	
	err = config.RedisClient.Set(ctx, userKey, userData, 0).Err() // No expiration
	if err != nil {
		return fmt.Errorf("failed to store user: %w", err)
	}
	
	// Store user session
	sessionKey := fmt.Sprintf("user:session:%s", user.PhoneNumber)
	sessionData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	
	// Session expires in 24 hours
	err = config.RedisClient.Set(ctx, sessionKey, sessionData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}
	
	return nil
}

// GetUser retrieves a user by phone number
func (r *RedisRepository) GetUser(ctx context.Context, phoneNumber string) (*domain.User, error) {
	key := fmt.Sprintf("user:profile:%s", phoneNumber)
	
	userData, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	var user domain.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	
	return &user, nil
}

// StoreReceipt stores a receipt in Redis
func (r *RedisRepository) StoreReceipt(ctx context.Context, receipt *domain.Receipt) error {
	// Store receipt details
	receiptKey := fmt.Sprintf("receipt:%s", receipt.ID)
	receiptData, err := json.Marshal(receipt)
	if err != nil {
		return fmt.Errorf("failed to marshal receipt: %w", err)
	}
	
	err = config.RedisClient.Set(ctx, receiptKey, receiptData, 0).Err() // No expiration
	if err != nil {
		return fmt.Errorf("failed to store receipt: %w", err)
	}
	
	// Add receipt ID to user's receipt list
	userReceiptsKey := fmt.Sprintf("user:receipts:%s", receipt.UserID)
	err = config.RedisClient.SAdd(ctx, userReceiptsKey, receipt.ID).Err()
	if err != nil {
		return fmt.Errorf("failed to add receipt to user list: %w", err)
	}
	
	return nil
}

// GetUserReceipts retrieves all receipts for a user
func (r *RedisRepository) GetUserReceipts(ctx context.Context, userID string) ([]*domain.Receipt, error) {
	userReceiptsKey := fmt.Sprintf("user:receipts:%s", userID)
	
	// Get all receipt IDs for the user
	receiptIDs, err := config.RedisClient.SMembers(ctx, userReceiptsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user receipt IDs: %w", err)
	}
	
	var receipts []*domain.Receipt
	
	// Fetch each receipt
	for _, receiptID := range receiptIDs {
		receiptKey := fmt.Sprintf("receipt:%s", receiptID)
		receiptData, err := config.RedisClient.Get(ctx, receiptKey).Result()
		if err != nil {
			continue // Skip if receipt not found
		}
		
		var receipt domain.Receipt
		err = json.Unmarshal([]byte(receiptData), &receipt)
		if err != nil {
			continue // Skip if receipt data is corrupted
		}
		
		receipts = append(receipts, &receipt)
	}
	
	return receipts, nil
}

