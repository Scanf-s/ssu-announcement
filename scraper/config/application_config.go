package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SSUAnnouncementURL string
	SSUPathURL         string
	DynamoDBClient     *dynamodb.Client
	DBTableName        string
	SSUPathID          string
	SSUPathPW          string
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("It is running in AWS Lambda.. Skipping")
	}

	// DynamoDB 클라이언트
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatal("DynamoDB config error : " + err.Error())
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	return &AppConfig{
		SSUAnnouncementURL: os.Getenv("SSU_ANNOUNCEMENT_URL"),
		SSUPathURL:         os.Getenv("SSU_PATH_URL"),
		DynamoDBClient:     dynamoClient,
		DBTableName:        os.Getenv("ANNOUNCEMENT_DB_NAME"),
		SSUPathID:          os.Getenv("SSU_PATH_ID"),
		SSUPathPW:          os.Getenv("SSU_PATH_PASSWORD"),
	}
}
