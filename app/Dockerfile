# Use the official Go image as the base image
FROM golang:latest AS build

# Set destination for COPY
WORKDIR /app

COPY . .

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/server.go

FROM alpine:latest as run


# Copy the application executable from the build image
COPY --from=build /app /app

WORKDIR /app

EXPOSE 3333

ENTRYPOINT ["./entrypoint.sh"]

CMD ["./main"]