.PHONY: server
.PHONY: client

clean:
	rm pb/*go

server:
	go run server/main.go

client:
	go run client/main.go

protoc:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

env:
	set -o allexport && source server/.env && set +o allexport

docker:
	cd server && docker-compose down && docker-compose up -d
	