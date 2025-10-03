package repository

import (
	"context"
	"log"
	"ssu-announcement-scraper/config"
	"ssu-announcement-scraper/internal/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func SaveScrapedData(ctx context.Context, cfg *config.AppConfig, data []dto.AnnouncementScrapedResult) {

	dbClient := cfg.DynamoDBClient
	tableName := cfg.DBTableName

	for _, item := range data {
		marshaledItem, err := attributevalue.MarshalMap(item)
		if err != nil {
			log.Println("Failed to marshal item:", err)
		}

		_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      marshaledItem,
		})
		if err != nil {
			log.Println("Failed to save item:", err)
		}
	}

	log.Println("Data saved to DynamoDB successfully")
}
