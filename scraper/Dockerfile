FROM golang:trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/scraper/main.go

FROM public.ecr.aws/lambda/provided:al2023 AS chromium_downloader
ARG TARGETARCH=amd64
# @sparticuz/chromium
# https://github.com/Sparticuz/chromium/releases
RUN dnf install -y tar && dnf clean all
ADD https://github.com/Sparticuz/chromium/releases/download/v141.0.0/chromium-v141.0.0-pack.x64.tar /tmp/chromium.tar
RUN tar -xf /tmp/chromium.tar -C /opt/

FROM public.ecr.aws/lambda/provided:al2023 AS runtime
WORKDIR ${LAMBDA_TASK_ROOT}
COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap
COPY --from=chromium_downloader /opt/ /opt/
CMD ["bootstrap"]