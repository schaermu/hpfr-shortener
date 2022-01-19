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
	make -j2 watch-go watch-svelte

watch-go:
	gotestsum --watch -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

watch-svelte:
	npm run --prefix ui test:watch

cover:
	gotestsum -f testname -- -tags=test ./... -coverprofile=coverage.out
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html