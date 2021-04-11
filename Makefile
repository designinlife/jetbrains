.PHONY: build clean release all
.DEFAULT_GOAL := all

OUTDIR="bin"
BINARY="jetbrains"
BINOUT="$OUTDIR/$BINARY"

all: clean build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINOUT} .

clean:
	@if [ -f ${BINOUT} ]; then rm ${BINOUT}; fi

release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINOUT} .
	# upx --brute ${BINOUT}
	upx bin/${BINOUT}
