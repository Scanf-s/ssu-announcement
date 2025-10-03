package repository

import (
	"context"
	"notifier/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetSubscribers(ctx context.Context, cfg *config.AppConfig, category string) ([]string, error) {
	// 구독자 목록가져오는 로직
	dbClient := cfg.DynamoDBClient
	tableName := cfg.DBTableName

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("CategoryIndex"),        // GSI 이름
		KeyConditionExpression: aws.String("Category = :category"), // 아래 ExpressionAttributeValues에서 매핑해줌
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category": &types.AttributeValueMemberS{Value: category}, // 카테고리 실제 값
		},
	}

	result, err := dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var emails []string
	for _, item := range result.Items { // 이메일 목록 추출
		email := item["Email"].(*types.AttributeValueMemberS).Value
		// *types.AttributeValueMemberS(tring) 타입으로 변환 후 Value 필드 접근해서 String 뽑아낼 수 있음
		emails = append(emails, email)
	}

	return emails, nil
}
