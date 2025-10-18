package config

import (
	"context"
	"log"
	"net/smtp"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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

func LoadConfig(ctx context.Context) *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// DynamoDB 클라이언트
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatal("DynamoDB config error : " + err.Error())
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// 환경변수 Setup
	var smtpHost, smtpPort, smtpUser, smtpPass, dbTableName string
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		ssmClient := ssm.NewFromConfig(cfg)
		parameterKeys := []string{
			"/smtp/host",
			"/smtp/port",
			"/smtp/user",
			"/smtp/password",
			"/ssu-announcement/subscribers-db-table-name",
		}
		resp, err := ssmClient.GetParameters(ctx, &ssm.GetParametersInput{
			Names:          parameterKeys,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatal("SSM Error : " + err.Error())
		}

		paramMap := make(map[string]string)
		for _, param := range resp.Parameters {
			paramMap[*param.Name] = *param.Value
		}

		smtpHost = paramMap["/smtp/host"]
		smtpPort = paramMap["/smtp/port"]
		smtpUser = paramMap["/smtp/user"]
		smtpPass = paramMap["/smtp/password"]
		dbTableName = paramMap["/ssu-announcement/subscribers-db-table-name"]
		if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || dbTableName == "" {
			log.Fatal("Required Parameters missing!")
		}
	} else {
		log.Fatal("This application cannot run locally!")
	}
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	return &AppConfig{
		SmtpHost:       smtpHost,
		SmtpPort:       smtpPort,
		SmtpUser:       smtpUser,
		SmtpPass:       smtpPass,
		Auth:           auth,
		DynamoDBClient: dynamoClient,
		DBTableName:    dbTableName, // 사용자 구독 정보 저장 테이블 이름
	}
}
