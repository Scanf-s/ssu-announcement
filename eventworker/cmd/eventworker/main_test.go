package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleRequest(t *testing.T) {
	// 환경 변수 설정
	os.Setenv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/test")
	os.Setenv("AWS_REGION", "ap-northeast-2")
	os.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
	os.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")

	// Mock DynamoDB Stream Event
	event := events.DynamoDBEvent{
		Records: []events.DynamoDBEventRecord{
			{
				EventName: "INSERT",
				Change: events.DynamoDBStreamRecord{
					NewImage: map[string]events.DynamoDBAttributeValue{
						"Link":       events.NewStringAttribute("https://example.com/1"),
						"Category":   events.NewStringAttribute("학사"),
						"Title":      events.NewStringAttribute("테스트 공지"),
						"Date":       events.NewStringAttribute("2025-10-04"),
						"Department": events.NewStringAttribute("교무처"),
						"Status":     events.NewStringAttribute("신규"),
					},
				},
			},
		},
	}

	result, err := handleRequest(context.Background(), event)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	t.Logf("Result: %s", result)
}
