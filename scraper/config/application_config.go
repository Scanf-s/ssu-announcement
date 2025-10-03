package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SSUAnnouncementURL string
	DynamoDBClient     *dynamodb.Client
	DBTableName        string
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

	var dynamoClient *dynamodb.Client
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		dynamoClient = dynamodb.NewFromConfig(cfg)
	} else { // 로컬 환경의 경우, 로컬 DynamoDB 엔드포인트 사용 (docker compose로 컨테이너 실행해주세요)
		dynamoClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://localhost:8000")
		})
	}

	return &AppConfig{
		SSUAnnouncementURL: os.Getenv("SSU_ANNOUNCEMENT_URL"),
		DynamoDBClient:     dynamoClient,
		DBTableName:        os.Getenv("DB_TABLE_NAME"),
	}
}
