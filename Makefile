.PHONY: build run clean

GOOS		?= linux
GOARCH		?= amd64

build:
		@GOOS=${GOOS} GOARCH=${GOARCH} go build -o crawler

run: build
		@./crawler

clean:
		@rm -f crawler || true
