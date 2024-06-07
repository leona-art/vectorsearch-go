FROM golang:1.22.4-alpine

WORKDIR /app

COPY backend/ ./

RUN go mod download

RUN go build -o /server ./cmd/server/main.go

EXPOSE 8080

CMD ["/server"]