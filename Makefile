# Eventworker 빌드
# -s -w 옵션 추가해서 디버깅 정보 제거하면 바이너리 크기가 줄어든다
build-eventworker:
	@echo "Building Eventworker service"
	cd eventworker && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap ./cmd/eventworker/main.go
	cd eventworker && zip eventworker.zip bootstrap
	@echo "Eventworker build completed"

# Clean
clean:
	@echo "Cleaning up"
	rm -f eventworker/bootstrap
	rm -f eventworker/eventworker.zip
	@echo "Cleanup completed"


# Eventworker 배포
deploy-eventworker:
	@echo "Deploying Eventworker service"
	aws lambda update-function-code \
	--function-name ${EVENTWORKER_LAMBDA_FUNCTION_NAME} \
	--zip-file fileb://eventworker/eventworker.zip \
	--region ${AWS_REGION}
	@echo "Eventworker deployment completed"