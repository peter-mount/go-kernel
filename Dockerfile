# Dockerfile used to test compile the libraries

# ============================================================
# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
#FROM docker.ceres.area51.dev/mirror/library/golang:alpine as build
FROM golang:alpine as build

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      git

# Static compile
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /work

# Seems go 1.16+ have changed things so this is required otherwise modules are not handled correctly with go.sum breaking
# via https://github.com/golang/go/issues/44129#issuecomment-860060061
RUN go env -w GOFLAGS=-mod=mod

# Download go module dependencies
COPY go.mod .
RUN go mod download

ADD . .

FROM build AS test

# Run each package separately so they don't interfere with each other
RUN CGO_ENABLED=0 go test -v .
RUN CGO_ENABLED=0 go test -v ./util
RUN CGO_ENABLED=0 go test -v ./test

# moduletest must be separate as it's single test will only work once
RUN CGO_ENABLED=0 go test -v ./test/moduletest
RUN CGO_ENABLED=0 go test -v ./test/interfaces

#FROM build AS compiler
#RUN CGO_ENABLED=0 go build -o test .
