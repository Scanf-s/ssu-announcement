FROM golang:trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/scraper/main.go

FROM public.ecr.aws/lambda/provided:al2023 AS runtime

# Chromium 설치
RUN dnf install -y epel-release && \
    dnf install -y chromium && \
    dnf clean all

# Chrome 심볼릭 링크 생성
RUN mkdir -p /opt/chrome && \
    ln -s /usr/bin/chromium-browser /opt/chrome/chrome

COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap
CMD ["bootstrap"]