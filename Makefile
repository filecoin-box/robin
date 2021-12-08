SHELL=/usr/bin/env bash

all: build-deps build
.PHONY: all build

unexport GOFLAGS

LOTUS_PATH:=./extern/lotus/

build-deps:
	git submodule update --init --recursive
	make -C ${LOTUS_PATH} deps

build: build-deps
	go mod tidy
	rm -rf robin
	go build -o robin ./main.go

.PHONY: clean
clean:
	-rm -f robin
