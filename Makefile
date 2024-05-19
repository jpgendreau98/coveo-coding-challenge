BINARY_NAME=tool


bin/$(BINARY_NAME): go.mod go.sum main.go cmd/* pkg/* 
	@go build -o bin/$(BINARY_NAME) main.go


build: bin/$(BINARY_NAME)
	@echo "build ok !"

build-windows:
	@env GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)_win main.go