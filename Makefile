generate-api:
	go install github.com/swaggo/swag/cmd/swag@v1.6.5
	swag init -g ./cmd/main.go -o docs

run-tests:
	go test -coverpkg=./... -coverprofile cover.out.tmp ./...
	cat cover.out.tmp | grep -v "/mock*" | grep -v "/cmd*" | grep -v "/docs*" | grep -v "/config*"> cover.out
	go tool cover -func cover.out

lint:
	golangci-lint run -c golangci.yml ./...
