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
			scrapedDataType := record.Change.NewImage["ScrapedDataType"].String()
			if scrapedDataType == "" {
				log.Println("ScrapedDataType is empty")
				continue
			}

			// Announcement인지 SSUPath인지 구분해서 적절한 DTO로 변환 후 SQS 전송
			var message dto.Message
			switch scrapedDataType {
			case "announcement":
				message = dto.AnnouncementMessage{
					Link:       record.Change.NewImage["Link"].String(),
					Category:   record.Change.NewImage["Category"].String(),
					Title:      record.Change.NewImage["Title"].String(),
					Date:       record.Change.NewImage["Date"].String(),
					Department: record.Change.NewImage["Department"].String(),
					Status:     record.Change.NewImage["Status"].String(),
				}
			case "ssu_path":
				message = dto.SSUPathMessage{
					Title:                   record.Change.NewImage["Title"].String(),
					Label:                   record.Change.NewImage["Label"].String(),
					Department:              record.Change.NewImage["Department"].String(),
					GradePointType:          record.Change.NewImage["GradePointType"].String(),
					Description:             record.Change.NewImage["Description"].String(),
					Image:                   record.Change.NewImage["Image"].String(),
					Link:                    record.Change.NewImage["Link"].String(),
					ApplicationPeriod:       record.Change.NewImage["ApplicationPeriod"].String(),
					EducationPeriod:         record.Change.NewImage["EducationPeriod"].String(),
					ApplicationTarget:       record.Change.NewImage["ApplicationTarget"].String(),
					ApplicationTargetStatus: record.Change.NewImage["ApplicationTargetStatus"].String(),
					Mileage:                 record.Change.NewImage["Mileage"].String(),
					Applicants:              record.Change.NewImage["Applicants"].String(),
					Waitlist:                record.Change.NewImage["Waitlist"].String(),
					Capacity:                record.Change.NewImage["Capacity"].String(),
				}
			default:
				log.Printf("Unknown ScrapedDataType: %s\n", scrapedDataType)
				continue
			}

			// SQS로 새로 추가된 데이터 전송 (notifier가 이거 받아서 알림 발송)
			err := service.SendMessageToSqs(ctx, cfg, message)
			if err != nil {
				log.Printf("Failed to send message to SQS: %s", err)
			}
		}
	}

	return "Success", nil
}

func main() {
	lambda.Start(handleRequest)
}
