run:
	@go run .

lint:
	@golangci-lint run ./...

tests:
	@go test -v ./test

cyclomatic:
	@gocyclo -over 7 .


clean:
	@rm -rf bin
