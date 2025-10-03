FROM golang:trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/ssu-announcement-notifier/main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap
ENTRYPOINT ["/lambda-entrypoint.sh"]