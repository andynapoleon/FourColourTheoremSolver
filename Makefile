build:
	@go build -o bin/goweb

run: build  
	@./bin/goweb

test:
	@go test -v ./...

mod:
	@go mod tidy 