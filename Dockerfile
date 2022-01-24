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
ADD . .

FROM build AS test

RUN CGO_ENABLED=0 go test -v . ./util ./test

#FROM build AS compiler
#RUN CGO_ENABLED=0 go build -o test .
