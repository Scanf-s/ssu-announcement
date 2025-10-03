package main

import (
	"context"
	"encoding/json"
	"log"
	"notifier/config"
	"notifier/internal/dto"
	"notifier/internal/repository"
	"notifier/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	cfg := config.LoadConfig()

	// SQS 이벤트 처리 로직
	for _, message := range sqsEvent.Records {
		// 메시지 Body를 Announcement로 파싱
		var announcement dto.Announcement
		err := json.Unmarshal([]byte(message.Body), &announcement)
		if err != nil {
			continue
		}

		// 구독자 조회
		emails, err := repository.GetSubscribers(ctx, cfg, announcement.Category)
		if err != nil {
			log.Printf("Error fetching subscribers: %v", err)
			continue
		}

		log.Printf("Category: %s, Title: %s, Subscribers: %v", announcement.Category, announcement.Title, emails)

		// 구독자들에게 이메일 알림 발송
		err = service.SendEmail(cfg, emails, announcement)
		if err != nil {
			log.Printf("Error sending emails: %v", err)
			continue
		}

		// 처리 로그 출력
		log.Printf("Processed message ID: %s", message.MessageId)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
