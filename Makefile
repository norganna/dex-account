all: dex-account

dex-account:
	@mkdir -p bin
	go build -o bin/$@ ./cmd

deps:
	go mod download

generate:
	go generate ./...

clean:
	rm bin/*

.PHONY: all
