FROM golang:trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/scraper/main.go

FROM public.ecr.aws/lambda/provided:al2023 AS runtime

# Chrome for Testing 설치
RUN dnf install -y -q findutils jq unzip && \
    dnf deplist https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm \
      | grep provider \
      | sort --unique \
      | awk '{print $2}' \
      | xargs dnf install --best --allowerasing --skip-broken -y -q >& /dev/null && \
    export CHROME_FOR_TESTING_RELEASE="$(curl --silent https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions-with-downloads.json | jq '.channels.Stable')" && \
    export CHROME_DOWNLOAD_URL="$(echo ""${CHROME_FOR_TESTING_RELEASE}"" | jq -r '.downloads.chrome[] | select(.platform == "linux64") | .url')" && \
    curl --remote-name --silent "${CHROME_DOWNLOAD_URL}" && \
    unzip -q chrome-linux64.zip && \
    ln --symbolic "${PWD}/chrome-linux64/chrome" /usr/local/bin/chrome && \
    export CHROMEDRIVER_DOWNLOAD_URL="$(echo ""${CHROME_FOR_TESTING_RELEASE}"" | jq -r '.downloads.chromedriver[] | select(.platform == "linux64") | .url')" && \
    curl --remote-name --silent "${CHROMEDRIVER_DOWNLOAD_URL}" && \
    unzip -q chromedriver-linux64.zip && \
    ln --symbolic "${PWD}/chromedriver-linux64/chromedriver" /usr/local/bin/chromedriver && \
    mkdir -p /opt/chrome && \
    ln --symbolic /usr/local/bin/chrome /opt/chrome/chrome && \
    dnf clean all

COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap
CMD ["bootstrap"]