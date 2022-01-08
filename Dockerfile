ARG GO_VERSION=1.17.6
ARG NODE_VERSION=16.13.1

# Step 1: build frontend
FROM node:${NODE_VERSION}-alpine as ui-builder
WORKDIR /src
COPY ui ./
RUN npm ci
RUN npm run build

# Step 2: build server
FROM golang:${GO_VERSION}-alpine as go-builder

# git is required to install using go mod
RUN apk update && apk add --no-cache git

WORKDIR /src

# download go dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy remainder (except ui) and compiled frontend from ui-builder
COPY [^ui]* .
COPY --from=ui-builder /src/dist ui/dist

# build static binary to run on distroless/static
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -installsuffix 'static' -o /app main.go

# Step 3: compile runtime image
FROM gcr.io/distroless/static

LABEL maintainer="schaermu"

USER nonroot:nonroot
COPY --from=go-builder --chown=nonroot:nonroot /app /app

ENTRYPOINT ["/app"]