#!/bin/bash

go mod download && \
go build -o ./build/img-sort

go mod download
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath \
      -ldflags="-w -s \
      -extldflags '-static'" -a \
      -o ./build/img-sort ./main.go