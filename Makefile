.PHONY: build clean release all
.DEFAULT_GOAL := all

BINARY="jetbrains"

all: clean build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/${BINARY} .

clean:
	@if [ -f ${BINARY} ]; then rm bin/${BINARY}; fi

release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/${BINARY} .
	# upx --brute ${BINARY}
	upx ${BINARY}
