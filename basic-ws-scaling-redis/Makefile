.PHONY: redis

run:
	go run main.go

redis:
	docker run -d -p 6379:6379 --name socket-cache redis:alpine