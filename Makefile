.PHONY: server

clean:
	rm pb/*go

server:
	go run server/main.go

protoc:
	protoc --go_out=. --go-grpc_out=. proto/*.proto
	