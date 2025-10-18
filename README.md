# 숭실대학교 공지사항 알림 서비스

카테고리별 구독 기반 실시간 이메일 알림 시스템
- SSU_PATH 동적 스크래퍼 미구현
- 프론트페이지 및 백엔드 API 미구현

## 📋 목차

- [프로젝트 개요](#프로젝트-개요)
- [아키텍처](#아키텍처)
- [기술스택](#기술스택)
- [구성 요소](#구성-요소)
- [환경 설정](#환경-설정)
- [배포](#배포)
- [유지보수](#유지보수)
- [트러블슈팅](#트러블슈팅)

---

## 프로젝트 개요

숭실대학교 공지사항을 스크래핑하여 카테고리별 구독자에게 이메일로 알림을 발송하는 서버리스 시스템입니다.

### 주요 기능

- 🔍 자동 공지사항 스크래핑 (3일 이내 데이터)
- 📧 카테고리별 구독 기반 이메일 알림
- ♻️ 중복 알림 방지
- ☁️ 서버리스 아키텍처

---

## 아키텍처
<img width="1321" height="780" alt="스크린샷 2025-10-04 오전 6 09 03" src="https://github.com/user-attachments/assets/befe7bff-39f2-472e-ac64-13115750dcea" />

---

## 기술스택
- **AWS Lambda:** 서버리스 함수 실행
- **AWS DynamoDB:** 공지사항 및 구독자 데이터 저장
- **AWS SQS:** 메시지 큐잉 서비스
- **AWS EventBridge:** 스케줄링 및 이벤트 라우팅
- **AWS IAM:** 권한 관리
- **AWS CloudWatch:** 로그 및 모니터링
- **Gmail SMTP:** 이메일 발송
- **Go:** 백엔드 로직 구현
- **GitHub Actions:** CI/CD 파이프라인
- **Docker:** 백엔드 이미지 빌드 및 배포
---

## 구성 요소

### 1. Scraper

**역할:** 숭실대학교 공지사항 페이지 스크래핑 → DynamoDB 저장

**주요 로직:**
- 3일 이내 공지사항 수집
- 중복 방지 (ConditionExpression: `attribute_not_exists(Link)`)
- EventBridge Schedule로 정기 실행

**환경 변수:**
```bash
SSU_ANNOUNCEMENT_URL=https://scatch.ssu.ac.kr/...
DYNAMODB_TABLE_NAME=데이터저장테이블
AWS_REGION=ap-northeast-2
```

**스크립트 위치:** `scraper/`
---

### 2. EventWorker

**역할:** DynamoDB Stream 이벤트 수신 → SQS로 메시지 전송

**주요 로직:**
- INSERT 이벤트만 처리 (신규 공지만)
- MessageBody: 전체 공지 정보 (JSON)
- MessageAttribute: Category (필터링용)

**환경 변수:**
```bash
SQS_QUEUE_URL=https://sqs.ap-northeast-2.amazonaws.com/ACCOUNT_ID/큐이름
```

**트리거 설정:**
- DynamoDB Stream: 데이터 저장되는 테이블 선택
- Batch size: 10

**스크립트 위치:** `eventworker/`
---

### 3. Notifier

**역할:** SQS 메시지 수신 -> 구독자 조회 -> 이메일 발송

**주요 로직:**
- SQS 메시지에서 공지 정보 파싱
- DynamoDB GSI로 카테고리별 구독자 조회
- Gmail SMTP로 이메일 발송

**환경 변수:**
```bash
DYNAMODB_TABLE_NAME=구독자관리테이블
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password  # Gmail 앱 비밀번호 (로그인할때 비밀번호가 아님)
AWS_REGION=ap-northeast-2
```

**트리거 설정:**
- SQS: 메세지 처리 큐 선택
- Batch size: 1 (중복 메일 전달 방지)
- Visibility timeout: 180초
- 동시 처리 수: 5(더 많이 필요하면 조정 가능)

**스크립트 위치:** `notifier/`
---

## 환경 설정

### 사전 준비

1. **AWS 계정** 및 IAM 사용자 생성 (귀찮으면 AdministratorAccess 권한 부여)
   - AmazonDynamoDBFullAccess
   - AmazonSQSFullAccess
   - AmazonEC2ContainerRegistryFullAccess
   - AWSLambda_FullAccess
   - AWSLambdaDynamoDBExecutionRole
   - AWSAPIGatewayAdministrator
   - 신뢰관계 설정은 아래 참고 or 알아서 선언 (브랜치, 레포지토리, AWS ARN 수정해야함)
    ```json 
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "Federated": "arn:aws:iam::AWS계정ARN:oidc-provider/token.actions.githubusercontent.com"
                },
                "Action": "sts:AssumeRoleWithWebIdentity",
                "Condition": {
                    "StringEquals": {
                        "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
                    },
                    "StringLike": {
                        "token.actions.githubusercontent.com:sub": [
                            "repo:깃허브계정이름/레포지토리이름:ref:refs/heads/notifier",
                            "repo:깃허브계정이름/레포지토리이름:ref:refs/heads/scraper",
                            "repo:깃허브계정이름/레포지토리이름:ref:refs/heads/eventworker"
                        ]
                    }
                }
            },
            {
                "Effect": "Allow",
                "Principal": {
                    "Federated": "arn:aws:iam::AWS계정ARN:oidc-provider/token.actions.githubusercontent.com"
                },
                "Action": "sts:AssumeRoleWithWebIdentity",
                "Condition": {
                    "StringEquals": {
                        "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
                    },
                    "StringLike": {
                        "token.actions.githubusercontent.com:sub": "repo:깃허브계정이름/레포지토리이름:ref:refs/heads/main"
                    }
                }
            },
            {
                "Effect": "Allow",
                "Principal": {
                    "Service": "lambda.amazonaws.com"
                },
                "Action": "sts:AssumeRole"
            }
        ]
    }
    ```
2. **AWS IAM에서 ID제공업체** 생성 (GitHub Actions 배포용 OIDC 생성)
3. **Gmail 앱 비밀번호** 생성
4. **ECR 레포지토리** 생성 (Docker 이미지 저장용)
5. **Github Actions** 돌려서 이미지 배포
   - `.github/workflows 참고
6. **AWS SAM CLI** 설치
7. 서버리스 리소스 배포
```bash
sam build
sam deploy --guided
```

---

## 백엔드 이미지 배포

### GitHub Actions 자동 배포

**워크플로우 위치:**
- `.github/workflows/scraper_workflow.yml`
- `.github/workflows/eventworker_workflow.yml`
- `.github/workflows/notifier_workflow.yml`

**배포 방법:**
```bash
# 각 모듈 이름으로 브랜치 만들어서 push
git add .
git commit -m "feat: update lambda function"
git push origin scraper/notifier/eventworker
```

**GitHub Secrets 설정 필요:**
```
AWS_ROLE_ARN=arn:aws:iam::ACCOUNT_ID:role/github-actions-role
SCRAPER_LAMBDA_FUNCTION_NAME=스크래퍼람다함수이름
EVENTWORKER_LAMBDA_FUNCTION_NAME=이벤트처리람다함수이름
NOTIFIER_LAMBDA_FUNCTION_NAME=알림처리람다함수이름
```
---

## 유지보수

### 1. 로그 확인

**CloudWatch Logs:**
- CloudWatch -> Logs -> 로그 그룹 선택 -> 각 람다 함수 이름으로 된 로그그룹 선택해서 로그 확인 가능

### 2. DynamoDB 데이터 확인

- AWS Console: DynamoDB → 테이블 선택 → 항목 탐색 -> 데이터베이스 테이블 선택 -> 스캔 or 쿼리해서 데이터 확인
---

## 트러블슈팅

### 문제 1: 이메일이 중복 발송됨

**원인:**
- Lambda timeout으로 인한 SQS 재시도

**해결:**
```bash
# 1. Lambda timeout 증가 (30초 이상)
aws lambda update-function-configuration \
    --function-name ssu-announcement-notifier \
    --timeout 30

# 2. SQS Visibility Timeout 조정 (Lambda timeout의 6배)
aws sqs set-queue-attributes \
    --queue-url YOUR_QUEUE_URL \
    --attributes VisibilityTimeout=180

# 3. 중복 메시지 제거
aws sqs purge-queue --queue-url YOUR_QUEUE_URL
```

### 문제 2: EventWorker가 호출 안됨

**원인:**
- DynamoDB Stream 트리거 미설정
- IAM 권한 부족

**해결:**
```bash
# 1. Stream 활성화 확인
aws dynamodb describe-table \
    --table-name ssu-announcement \
    --query 'Table.StreamSpecification'

# 2. Lambda 권한 확인 (AWSLambdaDynamoDBExecutionRole 필요)
aws iam list-attached-role-policies \
    --role-name eventworker-lambda-role
```

### 문제 3: 이메일 발송 실패 (SMTP 오류)

**원인:**
- Gmail 앱 비밀번호 오류
- 2단계 인증 미설정

**해결:**
1. Google 계정 → 보안 → 2단계 인증 활성화
2. 앱 비밀번호 생성: https://myaccount.google.com/apppasswords
3. Lambda 환경변수에 새 비밀번호 설정

### 문제 4: 구독자 조회 안됨

**원인:**
- GSI (CategoryIndex) 미생성
- 테이블 이름 불일치

**해결:**
```bash
# 1. GSI 확인
aws dynamodb describe-table \
    --table-name Subscriptions \
    --query 'Table.GlobalSecondaryIndexes'

# 2. 테스트 쿼리
aws dynamodb query \
    --table-name Subscriptions \
    --index-name CategoryIndex \
    --key-condition-expression "Category = :cat" \
    --expression-attribute-values '{":cat":{"S":"학사"}}'
```

### 문제 5: Scraper가 데이터를 못 가져옴

**원인:**
- 웹사이트 구조 변경
- URL 변경

**해결:**
1. 대상 웹사이트 직접 접속하여 구조 확인
2. `scraper/internal/scraper/ssu_announcement_scraper.go` 수정
3. `scraper/internal/service/ssu_announcement_parser/html_parser.go` 수정

---
