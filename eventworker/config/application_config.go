package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SqsClient *sqs.Client
	QueueUrl  string
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// SQS 클라이언트
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	sqsClient := sqs.NewFromConfig(cfg)

	return &AppConfig{
		SqsClient: sqsClient,
		QueueUrl:  os.Getenv("SQS_QUEUE_URL"),
	}
}
