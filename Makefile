build:
	@go build -o bin/hangman

run: build
	@./bin/hangman

test:
	@go test -v ./..