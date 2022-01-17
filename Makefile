generate-swagger-api:
	go install github.com/swaggo/swag/cmd/swag@v1.6.5
	swag init -g ./cmd/main.go -o docs

run-tests:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
