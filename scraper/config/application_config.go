package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SSUAnnouncementURL string
	SSUPathURL         string
	DynamoDBClient     *dynamodb.Client
	DBTableName        string
	SSUPathID          string
	SSUPathPW          string
	ChromeLauncher     *launcher.Launcher
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

	// Chrome 런처 설정
	var chromeLauncher *launcher.Launcher
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda 환경: Chrome 바이너리 경로 지정
		chromePath := "/opt/chrome/chrome" // 컨테이너 내부 파일시스템에 설치된 Chrome 경로
		chromeLauncher = launcher.New().
			Bin(chromePath).
			Headless(true).
			NoSandbox(true). // Lambda에서는 sandbox 모드 비활성화
			Set("disable-gpu").
			Set("disable-dev-shm-usage").
			Set("disable-setuid-sandbox").
			Set("no-first-run").
			Set("no-zygote").
			Set("single-process")
	} else {
		// 로컬 환경은 기본 설정 사용
		chromeLauncher = launcher.New().Headless(true)
	}

	return &AppConfig{
		SSUAnnouncementURL: os.Getenv("SSU_ANNOUNCEMENT_URL"),
		SSUPathURL:         os.Getenv("SSU_PATH_URL"),
		DynamoDBClient:     dynamoClient,
		DBTableName:        os.Getenv("ANNOUNCEMENT_DB_NAME"),
		SSUPathID:          os.Getenv("SSU_PATH_ID"),
		SSUPathPW:          os.Getenv("SSU_PATH_PASSWORD"),
		ChromeLauncher:     chromeLauncher,
	}
}
