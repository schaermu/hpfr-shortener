BINARY_NAME=hpfr-shortener

build:
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} main.go

run:
	./build/${BINARY_NAME}

start: build run

clean:
	go clean
	rm ./build

test:
	go test ./...

cover:
	go test ./... -coverprofile=coverage.out

docker:
	docker build .