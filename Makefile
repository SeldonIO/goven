lint:
	golangci-lint run --fix

test:
	go test ./...

test-coverage:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out