build:
	@go build -o bin/bankple

run: build
	@./bin/bankple

test:
	@go test -v ./...