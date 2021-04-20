.PHONY: build clean release linux windows all
.DEFAULT_GOAL := all

OUTDIR="bin"
BINARY="jetbrains"
BINOUT="${OUTDIR}/${BINARY}"

all: clean release linux windows

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINOUT} .

clean:
	@if [ -f ${BINOUT} ]; then rm ${BINOUT}; fi

release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINOUT} .
	# upx --brute ${BINOUT}
	upx ${BINOUT}

linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${OUTDIR}/${BINARY}-linux-amd64 .
	upx ${OUTDIR}/${BINARY}-linux-amd64

windows: clean
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${OUTDIR}/${BINARY}-windows-amd64.exe .
