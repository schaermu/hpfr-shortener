BINARY_NAME=hpfr-shortener

build:
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .

run:
	./build/${BINARY_NAME}

start: clean build run

clean:
	go clean
	rm -rf ./build

test:
	go test -tags=test ./...

cover:
	go test -tags=test ./... -coverprofile=coverage.out

docker:
	docker build .