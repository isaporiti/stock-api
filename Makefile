run:
	go run cmd/main.go

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out

docker_build:
	docker build -t stock-api .

docker_run:
	docker run --rm -p 5001:5001 stock-api