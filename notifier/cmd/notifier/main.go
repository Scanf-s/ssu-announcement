package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/mail.v2"
)

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	// 메일 설정
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USER")

	// SQS 이벤트 처리 로직
	for _, message := range sqsEvent.Records {
		// 메시지 속성 추출
		category := message.MessageAttributes["Category"].StringValue
		link := message.MessageAttributes["Link"].StringValue
		title := message.MessageAttributes["Title"].StringValue
		date := message.MessageAttributes["Date"].StringValue
		department := message.MessageAttributes["Department"].StringValue
		status := message.MessageAttributes["Status"].StringValue

		// 구독자 조회

		// 구독자들에게 이메일 알림 발송

		// 이벤트 메세지 삭제

	}
}

func main() {
	lambda.Start(handleRequest)
}
