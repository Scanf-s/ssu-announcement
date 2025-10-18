package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SqsClient *sqs.Client
	QueueUrl  string
}

func LoadConfig(ctx context.Context) *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// SQS 클라이언트
	sqsClient := sqs.NewFromConfig(cfg)

	// 환경변수 Setup
	var queueUrl string
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		ssmClient := ssm.NewFromConfig(cfg)
		parameterKeys := []string{
			"/sqs/queue-url",
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

		queueUrl = paramMap["/sqs/queue-url"]
		if queueUrl == "" {
			log.Fatal("Required Parameter missing!")
		}
	} else {
		log.Fatal("This script cannot run locally")
	}

	return &AppConfig{
		SqsClient: sqsClient,
		QueueUrl:  queueUrl,
	}
}
