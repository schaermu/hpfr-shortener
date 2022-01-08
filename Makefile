.PHONY: build

build:
	go build -o ./bin/hpfr-shortener ./

.PHONY: run

run: build
	./bin/hpfr-shortener