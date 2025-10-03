package main

import (
	"context"
	"eventworker/config"
	"eventworker/internal/dto"
	"eventworker/internal/service"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.DynamoDBEvent) (string, error) {
	log.Println("Event handler started")

	// 환경변수 로드
	cfg := config.LoadConfig()

	// DynamoDB 스트림에서 새롭게 추가된 항목들만 SQS로 전달
	for _, record := range event.Records {
		if record.EventName == "INSERT" && record.Change.NewImage != nil { // 새로운 항목이 추가된 경우
			message := dto.Message{
				Link:       record.Change.NewImage["Link"].String(),
				Category:   record.Change.NewImage["Category"].String(),
				Title:      record.Change.NewImage["Title"].String(),
				Date:       record.Change.NewImage["Date"].String(),
				Department: record.Change.NewImage["Department"].String(),
				Status:     record.Change.NewImage["Status"].String(),
			}

			// SQS로 새로 추가된 데이터 PK 전송 (notifier가 이거 받아서 알림 발송)
			err := service.SendMessageToSqs(ctx, cfg, message)
			if err != nil {
				log.Printf("Failed to send message to SQS: %s\n", err)
			}
		}
	}

	return "Success", nil
}

func main() {
	lambda.Start(handleRequest)
}
