build:
	@go build -o ./bin/spxcedrive

run: build
	@./bin/spxcedrive

test:
	@go test ./...