# Use the official Golang image from Docker Hub
FROM golang:latest

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Set the working directory
WORKDIR /app

# Copy the project files into the container
COPY . .

# Set the entry point to run Goose or keep it as default
CMD ["tail", "-f", "/dev/null"]
