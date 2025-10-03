package service

import (
	"context"
	"encoding/json"
	"eventworker/config"
	"eventworker/internal/dto"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func SendMessageToSqs(ctx context.Context, cfg *config.AppConfig, message dto.Message) error {
	// Message를 JSON으로 직렬화
	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// SQS로 새로 추가된 데이터 전송 (notifier가 이거 받아서 알림 발송)
	input := &sqs.SendMessageInput{
		QueueUrl:    &cfg.QueueUrl,
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Category": {
				DataType:    aws.String("String"),
				StringValue: aws.String(message.Category),
			},
		},
	}

	result, err := cfg.SqsClient.SendMessage(ctx, input)
	if err != nil {
		return err
	}
	log.Printf("Message sent to SQS: %s\n", *result.MessageId)
	return nil
}
