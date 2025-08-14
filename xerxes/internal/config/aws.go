package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var SNSClient *sns.Client

func InitAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	SNSClient = sns.NewFromConfig(cfg)
	log.Println("Successfully initialized AWS SNS client")
}

func SendSMS(phoneNumber, message string) error {
	ctx := context.Background()

	// Format phone number for international use
	// AWS SNS expects +1XXXXXXXXXX format for US numbers
	formattedPhone := formatPhoneNumber(phoneNumber)

	input := &sns.PublishInput{
		Message: &message,
		PhoneNumber: &formattedPhone,
	}

	_, err := SNSClient.Publish(ctx, input)
	return err
}

func formatPhoneNumber(phone string) string {
	clean := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			clean += string(char)
		}
	}

	// Add +1 prefix if it's a 10-digit US number
	if len(clean) == 10 {
		return "+1" + clean
	}
	
	// If it already has country code, just add +
	if len(clean) == 11 && clean[0] == '1' {
		return "+" + clean
	}

	// Otherwise, assume it's already formatted
	return "+" + clean
}