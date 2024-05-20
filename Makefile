build:
	@go build -o bin/wserv

run: build
	@./bin/wserv

test:
	@go test -v ./...
