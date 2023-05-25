.PHONY: garnet

build:
	@go build

run:
	@go build -o ./build/game ./cmd/game && ./build/game

run-p1:
	@go build -o ./build/game ./cmd/game && ./build/game user1 password1

run-p2:
	@go build -o ./build/game ./cmd/game && ./build/game user2 password2

run-indexer:
	@go build -o ./build/indexer ./cmd/indexer && ./build/indexer

run-localnet:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 && cd contracts-builder/contracts && pnpm run devnode

contracts:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 && cd contracts-builder/contracts && pnpm run dev && cd ../.. && cp contracts-builder/contracts/out/IWorld.sol/IWorld.abi.json internal/txbuilder/


run-generator:
	@go build -o ./build/generator ./cmd/generator && ./build/generator

init-contracts:
	@source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 && cd contracts-builder/contracts && pnpm run install

