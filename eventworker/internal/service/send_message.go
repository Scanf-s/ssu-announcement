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

func SendMessageToSqs[T dto.Message](ctx context.Context, cfg *config.AppConfig, message T) error {

	// Message를 JSON으로 직렬화
	messageBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to serialize message: %v\n", err)
		return err
	}

	// SQS로 새로 추가된 데이터 전송 (notifier가 이거 받아서 알림 발송)
	// Notifier에서 json.Unmarshal(message.Body, &announcement or &ssu_path)로 사용해주면 됨
	var messageAttributes map[string]types.MessageAttributeValue
	switch messageType := message.GetMessageType(); messageType {
	case dto.MessageTypeAnnouncement:
		// AnnouncementMessage로 타입 설정
		announcementMsg, ok := any(message).(dto.AnnouncementMessage)
		if !ok {
			log.Println("failed to cast message to AnnouncementMessage")
			return nil
		}
		// AnnouncementMessage의 경우 Category 속성으로 구독자별로 이메일 전송하므로 필요
		messageAttributes = map[string]types.MessageAttributeValue{
			"MessageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(dto.MessageTypeAnnouncement)),
			},
			"Category": {
				DataType:    aws.String("String"),
				StringValue: aws.String(announcementMsg.Category),
			},
		}
	case dto.MessageTypeSSUPath:
		// SSUPathMessage로 타입 설정
		_, ok := any(message).(dto.SSUPathMessage)
		if !ok {
			log.Println("failed to cast message to SSUPathMessage")
			return nil
		}
		messageAttributes = map[string]types.MessageAttributeValue{
			"MessageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(dto.MessageTypeSSUPath)),
			},
		}
	default:
		log.Printf("unsupported message type: %s", messageType)
		return nil
	}

	input := &sqs.SendMessageInput{
		QueueUrl:          &cfg.QueueUrl,
		MessageBody:       aws.String(string(messageBody)),
		MessageAttributes: messageAttributes,
	}

	result, err := cfg.SqsClient.SendMessage(ctx, input)
	if err != nil {
		log.Printf("Failed to send message to SQS: %v\n", err)
		return err
	}

	log.Printf("Message sent to SQS: %s", *result.MessageId)
	return nil
}
