build:
	@go build -buildvcs=false -o bin/wserv

run: build
	@./bin/wserv

test:
	@go test -v ./...
