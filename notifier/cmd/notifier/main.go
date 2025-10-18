package main

import (
	"context"
	"log"
	"notifier/config"
	"notifier/internal/dto"
	"notifier/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	cfg := config.LoadConfig(ctx)

	// SQS 이벤트 처리 로직
	for _, message := range sqsEvent.Records {
		if message.MessageAttributes["MessageType"].StringValue != nil {
			var messageType = *message.MessageAttributes["MessageType"].StringValue
			switch messageType {
			case string(dto.MessageTypeAnnouncement):
				err := service.NotificationService(ctx, cfg, message.Body, *message.MessageAttributes["Category"].StringValue)
				if err != nil {
					log.Printf("Error sending email: %v", err)
				}
			case string(dto.MessageTypeSSUPath):
				err := service.NotificationService(ctx, cfg, message.Body, "ssu_path")
				if err != nil {
					log.Printf("Error sending email: %v", err)
				}
			default:
				log.Printf("Unsupported message type: %s", messageType)
			}

			// 처리 로그 출력
			log.Printf("Processed message ID: %s", message.MessageId)
		} else {
			log.Printf("Invalid Message Attributes")
		}
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
