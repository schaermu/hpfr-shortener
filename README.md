[![CircleCI](https://circleci.com/gh/schaermu/hpfr-shortener/tree/main.svg?style=shield)](https://circleci.com/gh/schaermu/hpfr-shortener/tree/main)
[![codecov](https://codecov.io/gh/schaermu/hpfr-shortener/branch/main/graph/badge.svg?token=QC1WL6JQTQ)](https://codecov.io/gh/schaermu/hpfr-shortener)
# hpfr-shortener
Another link shortener. It isn't a problem that deserves solving and it sure is not rocket science, but thats not the purpose. This is simply a small project used by me to learn a bit more about [Svelte](https://svelte.dev/) and [Golang](https://go.dev/) as a web platform.

# Setup
1. Clone the repository: `git clone https://github.com/schaermu/hpfr-shortener.git`
2. (optional) Start a local MongoDB server: `docker-compose up -d`
2. Build and run the app: `make start`

# Build
Run `docker build .` to build a production-ready, [distroless](https://github.com/GoogleContainerTools/distroless)-based docker image. The latest version of the application is also published to the [Github Registry](https://github.com/schaermu/hpfr-shortener/pkgs/container/hpfr-shortener).

# Testing
## Single run
* Run all tests: `make test`
* Only run Svelte tests: `make test-svelte`
* Only run Go tests: `make test-go`

## Watch mode
* Watch all tests: `make watch`
* Only watch Svelte tests: `make watch-svelte`
* Only watch Go tests: `make watch-go`