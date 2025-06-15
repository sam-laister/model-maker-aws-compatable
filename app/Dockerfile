FROM golang:alpine3.21 AS builder

RUN apk update && apk add --no-cache \
    git \
    build-base \
    jq \
    cmake 

RUN go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go install github.com/joho/godotenv/cmd/godotenv@latest && \
    go install github.com/nikolaydubina/go-cover-treemap@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

EXPOSE 3333

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]

RUN make build

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3333/health || exit 1

CMD ["make", "run"]