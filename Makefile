.PHONY: garnet

build:
	@go build

run:
	@go build -o ./build/game ./cmd/game && ./build/game

run-indexer:
	@go build -o ./build/indexer ./cmd/indexer && ./build/indexer

run-localnet:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 && cd contracts-builder/contracts && pnpm run devnode

contracts:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 && cd contracts-builder/contracts && pnpm run dev

init-contracts:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 && cd contracts-builder/contracts && pnpm run install

