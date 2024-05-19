BINARY_NAME=tool


bin/$(BINARY_NAME): go.mod go.sum main.go cmd/* pkg/** pkg/aws/* pkg/util/*
	@go build -o bin/$(BINARY_NAME) main.go

go.mod:
	@go mod init projet-devops-coveo/$(BINARY_NAME)


tests: coverage.txt
	@go tool cover -func coverage.txt
	@rm -rf coverage.txt


go.sum: go.mod
	@go mod tidy


build: bin/$(BINARY_NAME)
	@echo "build ok !"

coverage.txt: go.sum
	@go test -coverprofile coverage.txt.tmp ./...
	@cat coverage.txt.tmp | grep -v "_mocks.go" > coverage.txt
	@rm -fr coverage.txt.tmp*


build-windows:
	@env GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)_win main.go

	tests-html: coverage.txt
	@go tool cover -html=coverage.txt -o ./reports/coverage.html
	@rm -rf coverage.txt
