FROM golang:1.23.2

# Set environment variables
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct

# Install swag tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Create a directory for the app
WORKDIR /app

# Copy project files (only if needed, else it's omitted)
COPY . .

ENTRYPOINT ["swag"]