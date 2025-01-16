FROM golang:1.23.2

# Set environment variables
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct

# Install golang-migrate tool with MySQL support
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create a directory for the app
WORKDIR /app

# Copy project files (only if needed, else it's omitted)
COPY . .

ENTRYPOINT ["migrate"]
