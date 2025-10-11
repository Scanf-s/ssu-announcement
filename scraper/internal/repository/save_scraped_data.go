package repository

import (
	"context"
	"log"
	"scraper/config"
	"scraper/internal/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func SaveScrapedData[T dto.ScrapedResult](ctx context.Context, cfg *config.AppConfig, data []T) error {

	dbClient := cfg.DynamoDBClient
	tableName := cfg.DBTableName

	// 데이터 저장
	// BatchWriteItem으로 한번에 25개씩 묶어서 처리할 수 있는데, 어짜피 첫페이지만 긁어오는거라서 25개가 안됨 -> 그냥 PutItem으로 하나씩 넣어주기
	// DynamoDB의 Link 속성이 PK입니다 (template.yaml 확인)
	for _, item := range data {
		marshaledItem, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}

		_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:           aws.String(tableName),
			Item:                marshaledItem,
			ConditionExpression: aws.String("attribute_not_exists(Link)"), // 해당 공지가 디비에 없을 때만 추가
		})
		if err != nil {
			return err
		}
	}

	log.Println("Data saved to DynamoDB successfully")
	return nil
}
