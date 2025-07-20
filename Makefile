.PHONY: build run clean

build:
		@go build -o crawler

run: build
		@./crawler

clean:
		rm -f crawler || true
