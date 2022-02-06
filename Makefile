BINARY_NAME=hpfr-shortener
GOCOVER=go tool cover

.PHONY: all
build:
	npm run --prefix ui build
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .

run:
	./build/${BINARY_NAME}

start: clean build run

clean:
	go clean
	go clean -testcache
	rm -rf ./build
	rm -rf ./ui/dist

watch-go:
	gow run .

watch-web:
	npm run --prefix ui dev

test: test-go test-web

test-go:
	gotestsum -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

test-web:
	npm run --prefix ui test

test-watch:
	$(MAKE) -j2 test-watch-go test-watch-web

test-watch-go:
	gotestsum --watch -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

test-watch-web:
	npm run --prefix ui test:watch
