package config

import (
	"context"
	"log"
	"net/smtp"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SmtpHost string
	SmtpPort string
	Auth     smtp.Auth
	SmtpUser string
	SmtpPass string

	DynamoDBClient *dynamodb.Client
	DBTableName    string
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// DynamoDB 클라이언트
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatal("DynamoDB config error : " + err.Error())
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	return &AppConfig{
		DynamoDBClient: dynamoClient,
		DBTableName:    os.Getenv("SUBSCRIBE_DB_TABLE_NAME"), // 사용자 구독 정보 저장 테이블 이름
	}
}
