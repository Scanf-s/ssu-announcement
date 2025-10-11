package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/aws/aws-lambda-go/lambda"

	"scraper/config"
	"scraper/internal/scraper"
)

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(handleRequest)
	} else {
		// 로컬 환경에서 실행할 때는 빈 context 전달
		_, err := handleRequest(context.Background(), nil)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func handleRequest(ctx context.Context, event json.RawMessage) (string, error) {
	log.Println("Running in Lambda environment")

	// 환경변수 불러오기
	cfg := config.LoadConfig(ctx)

	// errgroup 생성
	g, gCtx := errgroup.WithContext(ctx)

	// 숭실대학교 공지사항 스크래핑 작업
	g.Go(func() error {
		return scraper.ScrapeSSUAnnouncements(gCtx, cfg)
	})

	// SSU-Path 프로그램 스크래핑 작업
	g.Go(func() error {
		return scraper.ScrapeSSUPathPrograms(gCtx, cfg)
	})

	// 모든 작업 완료 대기 (하나라도 실패하면 에러 반환)
	if err := g.Wait(); err != nil {
		log.Printf("Scraping failed: %v", err)
		return "", err
	}

	log.Println("All scraping tasks completed successfully")
	return "Success", nil
}
