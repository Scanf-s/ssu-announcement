FROM golang:trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/scraper/main.go

FROM public.ecr.aws/lambda/provided:al2023 AS runtime

RUN dnf -y update && \
    dnf install -y \
    wget \
    unzip \
    nss \
    atk \
    cups-libs \
    libdrm \
    libxkbcommon \
    libXcomposite \
    libXcursor \
    libXdamage \
    libXext \
    libXi \
    libXScrnSaver \
    libXtst \
    pango \
    xorg-x11-server-Xvfb \
    gtk3 \
    libxshmfence \
    alsa-lib \
    ipa-gothic-fonts \
    xdg-utils \
    liberation-fonts > /dev/null 2>&1

# Chrome 설치
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm && \
    rpm install ./google-chrome-stable_current_x86_64.rpm || true && \
    mkdir -p /opt/chrome && \
    ln -s /usr/bin/google-chrome /opt/chrome/chrome

COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap
CMD ["bootstrap"]