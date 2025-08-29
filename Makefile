.PHONY: build run clean

BIN := bin/hexlet-path-size$(EXE)
EXE :=

ifeq ($(OS),Windows_NT)
    EXE := .exe
endif

build: $(BIN)

$(BIN):
	mkdir -p bin
	go build -o $(BIN) ./cmd/hexlet-path-size

run: build
	$(BIN)

clean:
	rm -rf bin