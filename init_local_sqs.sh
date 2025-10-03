#!/bin/sh

echo "Waiting for LocalStack to be ready..."
sleep 5

echo "Creating SQS queue..."
aws sqs create-queue \
    --queue-name test \
    --endpoint-url http://localstack:4566 \
    --region ap-northeast-2

echo "SQS queue created successfully!"
