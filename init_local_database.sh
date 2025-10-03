#!/bin/sh

echo "Waiting for DynamoDB Local to be ready..."
sleep 3

echo "Creating DynamoDB table..."
aws dynamodb create-table \
    --table-name test \
    --attribute-definitions AttributeName=Link,AttributeType=S \
    --key-schema AttributeName=Link,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://dynamodb-local:8000

echo "Table created successfully!"