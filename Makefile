BINARY_NAME=hpfr-shortener

GOCOVER=go tool cover

build:
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .

run:
	./build/${BINARY_NAME}

start: clean build run

clean:
	go clean
	go clean -testcache
	rm -rf ./build

test:
	gotestsum -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

watch:
	gotestsum --watch -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

cover:
	gotestsum -f testname -- -tags=test ./... -coverprofile=coverage.out
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html