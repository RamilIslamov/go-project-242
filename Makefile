# Makefile
.PHONY: build run clean
BIN := bin/hexlet-path-size.exe

build: $(BIN)

$(BIN):
	@if not exist bin mkdir bin
	go build -o $(BIN) ./cmd/hexlet-path-size

run: build
	$(BIN)

clean:
	@if exist bin rmdir /S /Q bin