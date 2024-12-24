# Use the official Go image as the base image
FROM golang:latest AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest as run

# Copy the application executable from the build image
COPY --from=build /main /main

WORKDIR /app
EXPOSE 3333
CMD ["/main"]