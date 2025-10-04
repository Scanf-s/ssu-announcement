package handler

import (
	"context"
	"encoding/json"
	"errors"

	"api/config"

	"api/internal/repository"

	"github.com/aws/aws-lambda-go/events"
)

func HandleGetSubscribes(ctx context.Context, cfg *config.AppConfig, request events.APIGatewayProxyRequest) ([]byte, error) {

	// Cognito에서 JWT발급할 때 sub에 사용자 이메일 넣어놓았음
	userEmail := request.RequestContext.Authorizer["claims"].(map[string]interface{})["sub"].(string)

	// DynamoDB에서 해당 사용자의 구독 목록 조회
	subscribes := repository.GetSubscribes(ctx, cfg, userEmail)

	// 구독 목록을 JSON 형식으로 반환
	response, err := json.Marshal(subscribes)
	if err != nil {
		return nil, errors.New("failed to serialized subscribed data")
	}
	return response, nil
}
