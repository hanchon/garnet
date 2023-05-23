# Garnet

Autonomous Worlds hackathon entry.

- Garnet: is a custom indexer and backend code running on top of MUDv2.
- Eternal Legends is the example game:
  - Eternal Legends, a blockchain-based turn-based tactical game, combines MUD2, Solidity, and GoLang. Players select 3 heroes to protect their castles in battles fought on an ASCII graphical interface. Destroying the opponent's base leads to victory.

## Requirements

- Go 1.19

To execute the MUD's scripts it's also required:

- Foundry toolkit
- Pnpm

The `Makefile` assumes that you have `nvm` as your node version handler and that it's installed at `/opt/homebrew/opt/nvm/nvm.sh`.

If you have `pnpm` installed globally you can remove `source /opt/homebrew/opt/nvm/nvm.sh && nvm use v18.12.0 &&` from the `Makefile`.

- Clone and install dependencies:

```sh
git clone https://github.com/hanchon/garnet
make init-contracts
```

## How to run it

- We need to run a local blockchain to deploy the contracts:
  - `make run-localnet`. It runs a local blockchain using `foundry`
- Compile and deploy your contracts:
  - `make contracts`. It runs the `mud-cli` in the contracts folder (`contracts-builder/contracts`), deploys them to the local blockchain, and copies the `IWorld.abi.json` file into the `Go` codebase.
- Run the indexer/server:
  - `make run-indexer`. It connects to the local blockchain and requests all the events related to MUD's transactions, parses those events and registers everything in an in-memory database.
    - The indexer also hosts a WebSocket server so clients can connect to it.
    - The WebSocket can create and send transactions to the blockchain.
- Run the client:
  - `make run`. It runs the client, it expects the username and password as arguments:
    - `make run-p1`. It runs the client but with the user1 credentials already set.
    - `make run-p2`. It runs the client but with the user2 credentials already set.
    - The client gets updates using a WebSocket, when it needs to interact with the blockchain it just sends a message to the WebSocket and the backend will generate and sign the transaction for the current user.

## Future developments

### Game

- Add unit boosters (like hp and dmg) that will spawn in random places every X turns.
- Add one special active ability to each unit (dmg).
- Add one passive ability to each unit (like 1 free movement each turn).

### Backend

- Add tests for the indexer.
- Add support for MUD's ephemeral events.
- Handle all the WebSocket possible errors.
- Instead of hardcoding the WorldID, it should be a parameter.

### Client

- Create a Unity client to connect to the backend and display the information with better graphics.

### Core

- Users just use their wallets to `create` and `join` matches.

- `Create` and `Join` matches transactions will delegate power to an `admin wallet` to send transactions representing each user (this wallet will be unique to each Match and controlled by the backend).

- Instead of reading events from the blocks, we are going to predict the events using the `Mempool` transactions.

- The admin wallet will send all the transactions related to a Match so the order will be deterministic (it will be controlled by the wallet's nonce).

- We control the wallet so we can resend transactions in case something fails to execute when included in a block.

- Why? Reading from the `Mempool` will move all the predictions from the client code to the backend code, and the predictions will be broadcasted to all the clients connected instead of just being executed to the player that sent the transaction.
