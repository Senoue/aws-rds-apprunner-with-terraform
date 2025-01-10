FROM golang:1.23.2

WORKDIR /app

# Enable Go modules
ENV GO111MODULE=on

# Install air
RUN go install github.com/air-verse/air@latest 
ENV PATH="/go/bin:${PATH}"

CMD ["air", "-c", ".air.toml"]