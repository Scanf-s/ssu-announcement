package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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

func LoadConfig(ctx context.Context) *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("It is running in AWS Lambda.. Skipping")
	}

	// DynamoDB 클라이언트
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatal("DynamoDB config error : " + err.Error())
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// 환경 변수 가져오기
	var ssuAnnouncementURL, ssuPathURL, dbTableName, ssuPathId, ssuPathPassword string

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda 환경: SSM Parameter Store에서 모두 가져오기
		ssmClient := ssm.NewFromConfig(cfg)
		parameterKeys := []string{
			"/ssu-announcement/url",
			"/ssu-announcement/db-table-name",
			"/ssu-path/url",
			"/ssu-path/student-id",
			"/ssu-path/password",
		}
		resp, err := ssmClient.GetParameters(ctx, &ssm.GetParametersInput{
			Names:          parameterKeys,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatal("SSM Error : " + err.Error())
		}

		// Parameter 매핑
		paramMap := make(map[string]string)
		for _, param := range resp.Parameters {
			paramMap[*param.Name] = *param.Value
		}

		ssuAnnouncementURL = paramMap["/ssu-announcement/url"]
		dbTableName = paramMap["/ssu-announcement/db-table-name"]
		ssuPathURL = paramMap["/ssu-path/url"]
		ssuPathId = paramMap["/ssu-path/student-id"]
		ssuPathPassword = paramMap["/ssu-path/password"]

		// 필수 Parameter 체크
		if ssuAnnouncementURL == "" || dbTableName == "" || ssuPathURL == "" || ssuPathId == "" || ssuPathPassword == "" {
			log.Fatal("Required SSM Parameters not found. Please check Parameter Store.")
		}
	} else {
		// 로컬 환경은 .env 파일에서 가져오기
		ssuAnnouncementURL = os.Getenv("SSU_ANNOUNCEMENT_URL")
		ssuPathURL = os.Getenv("SSU_PATH_URL")
		dbTableName = os.Getenv("ANNOUNCEMENT_DB_NAME")
		ssuPathId = os.Getenv("SSU_PATH_ID")
		ssuPathPassword = os.Getenv("SSU_PATH_PASSWORD")

		if ssuAnnouncementURL == "" || ssuPathURL == "" || dbTableName == "" || ssuPathId == "" || ssuPathPassword == "" {
			log.Fatal("Required environment variables not set in .env file")
		}
	}

	// Chrome 런처 설정
	var chromeLauncher *launcher.Launcher
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda 환경: Chrome 바이너리 경로 지정
		chromePath := "/opt/chromium" // 컨테이너 내부 파일시스템에 설치된 Chrome 경로
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
		SSUAnnouncementURL: ssuAnnouncementURL,
		SSUPathURL:         ssuPathURL,
		DynamoDBClient:     dynamoClient,
		DBTableName:        dbTableName,
		SSUPathID:          ssuPathId,
		SSUPathPW:          ssuPathPassword,
		ChromeLauncher:     chromeLauncher,
	}
}
