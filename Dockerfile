FROM registry.hub.docker.com/library/golang:latest

RUN apt-get update && apt-get install -y \
    curl \
    iputils-ping \
    telnet \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

EXPOSE 8080

CMD ["./cmd/main"]