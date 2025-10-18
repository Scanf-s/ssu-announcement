EVENTWORKER_LAMBDA_FUNCTION_NAME ?= asdf
AWS_REGION ?= ap-northeast-2

# Eventworker 빌드
# -s -w 옵션 추가해서 디버깅 정보 제거하면 바이너리 크기가 줄어든다
build-eventworker:
	@echo "Building Eventworker service"
	cd eventworker && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap ./cmd/eventworker/main.go
	cd eventworker && zip eventworker.zip bootstrap
	@echo "Eventworker build completed"

# Eventworker 배포
deploy-eventworker:
	@echo "Deploying Eventworker service"
	aws lambda update-function-code \
	--function-name $(EVENTWORKER_LAMBDA_FUNCTION_NAME) \
	--zip-file fileb://eventworker/eventworker.zip \
	--region $(AWS_REGION)
	@echo "Eventworker deployment completed"

# Eventworker 빌드 정리
clean-eventworker:
	@echo "Cleaning up"
	rm -f eventworker/bootstrap
	rm -f eventworker/eventworker.zip
	@echo "Cleanup completed"

# Notifier 빌드
build-notifier:
	@echo "Building Notifier service"
	cd notifier && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap ./cmd/notifier/main.go
	cd notifier && zip notifier.zip bootstrap
	@echo "Notifier build completed"

# Notifier 배포
deploy-notifier:
	@echo "Deploying Notifier service"
	aws lambda update-function-code \
	--function-name $(NOTIFIER_LAMBDA_FUNCTION_NAME) \
	--zip-file fileb://notifier/notifier.zip \
	--region $(AWS_REGION)
	@echo "Notifier deployment completed"

# Notifier 빌드 정리
clean-notifier:
	@echo "Cleaning up"
	rm -f notifier/bootstrap
	rm -f notifier/notifier.zip
	@echo "Cleanup completed"