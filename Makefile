.PHONY: garnet

build:
	@go build

run:
	@go build -o ./build/game ./cmd/game && ./build/game

run-indexer:
	@go build -o ./build/indexer ./cmd/indexer && ./build/indexer

