package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"scraper/config"
	"scraper/internal/repository"
	"scraper/internal/scraper"
	"scraper/internal/service/ssu_announcement_parser"
)

func main() {
	// AWS_LAMBDA_FUNCTION_NAME 환경변수로 Lambda 환경 판별
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda 환경
		lambda.Start(handleRequest)
	} else {
		// 로컬 환경
		runLocal()
	}
}

func handleRequest(ctx context.Context, event json.RawMessage) (string, error) {
	log.Println("Running in Lambda environment")

	// 환경변수 불러오기
	cfg := config.LoadConfig()

	// 숭실대학교 공지사항 스크래핑
	resultHtml := scraper.ScrapeSSUAnnouncements(cfg)

	// 공지사항 HTML 파싱해서 원하는 정보 추출
	parsedResult := ssu_announcement_parser.ParseSSUAnnouncementsHtml(resultHtml)
	log.Println(parsedResult)

	// DynamoDB에 저장
	repository.SaveScrapedData(ctx, cfg, parsedResult)

	return "Success", nil
}

func runLocal() {
	log.Println("Running in local environment")

	// 환경변수 불러오기
	cfg := config.LoadConfig()

	// 숭실대학교 공지사항 스크래핑
	resultHtml := scraper.ScrapeSSUAnnouncements(cfg)

	// 공지사항 HTML 파싱해서 원하는 정보 추출
	parsedResult := ssu_announcement_parser.ParseSSUAnnouncementsHtml(resultHtml)
	log.Println(parsedResult)

	// DynamoDB에 저장
	repository.SaveScrapedData(context.TODO(), cfg, parsedResult)
}
