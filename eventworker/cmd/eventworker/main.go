package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func handleRequest(ctx context.Context, event events.DynamoDBEvent) (string, error) {
	log.Println("Running in Lambda environment")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	queueURL := os.Getenv("SQS_QUEUE_URL")

	// DynamoDB 스트림에서 새롭게 추가된 항목들만 SQS로 전달
	for _, record := range event.Records {
		if record.EventName == "INSERT" && record.Change.NewImage != nil { // 새로운 항목이 추가된 경우
			link := record.Change.NewImage["Link"].String()
			category := record.Change.NewImage["Category"].String()

			// SQS로 새로 추가된 데이터 PK 전송 (notifier가 이거 받아서 db 조회 후 알림 발송)
			input := &sqs.SendMessageInput{
				QueueUrl:    &queueURL,
				MessageBody: aws.String(link),
				MessageAttributes: map[string]types.MessageAttributeValue{ // 카테고리 정보 전송 (카테고리 기준 구독 필터링)
					"Category": {DataType: aws.String("String"), StringValue: aws.String(category)}},
			}

			result, err := sqsClient.SendMessage(ctx, input)
			if err != nil {
				log.Printf("Failed to send message to SQS: %s\n", err)
				continue
			}
			log.Printf("Message sent to SQS: %s\n", *result.MessageId)
		}
	}

	return "Success", nil
}

func main() {
	lambda.Start(handleRequest)
}
