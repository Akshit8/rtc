# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o chat-app main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/chat-app .
EXPOSE 8080
CMD ["./chat-app"]