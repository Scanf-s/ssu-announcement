package main

import (
	"api/config"
	"context"
	"log"

	"api/internal/handler"

	"github.com/aws/aws-lambda-go/events"
)

func main() {
	log.Println("SSU Announcement Service API")
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("API handler started from %s", request.RequestContext.Identity.SourceIP)
	cfg := config.LoadConfig()

	// 라우팅
	switch {
	case request.HTTPMethod == "GET" && request.Path == "/subscribes":
		// 본인의 구독 목록 조회
		// Cognito에서 Authorizer에다가 사용자 정보 넣어주니까 그거 사용해서 반환
		return handler.HandleGetSubscribes(ctx, cfg, request)
	case request.HTTPMethod == "POST" && request.Path == "/subscribes":
		// 구독 추가
		return handler.HandleCreateSubscribe(ctx, cfg, request)
	case request.HTTPMethod == "PATCH" && request.Path == "/subscribes":
		// 구독 수정
		return handler.HandleUpdateSubscribe(ctx, cfg, request)
	case request.HTTPMethod == "DELETE" && request.Path == "/subscribes":
		// 구독 삭제
		return handler.HandleDeleteSubscribe(ctx, cfg, request)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       `{"error": "Not Found"}`,
	}, nil
}
