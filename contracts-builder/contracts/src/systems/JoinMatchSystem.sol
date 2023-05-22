// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import {System} from "@latticexyz/world/src/System.sol";
import {Match} from "../codegen/tables/Match.sol";
import {PlayerOne} from "../codegen/tables/PlayerOne.sol";
import {PlayerTwo} from "../codegen/tables/PlayerTwo.sol";
import {CurrentTurn} from "../codegen/tables/CurrentTurn.sol";
import {CurrentPlayer} from "../codegen/tables/CurrentPlayer.sol";
import {CurrentMana} from "../codegen/tables/CurrentMana.sol";
import {addressToEntityKey} from "../addressToEntityKey.sol";
import {positionToEntityKey} from "../positionToEntityKey.sol";
import {CardTypes} from "../codegen/Types.sol";
// Tables
import {Card} from "../codegen/tables/Card.sol";
import {OwnedBy} from "../codegen/tables/OwnedBy.sol";
import {UsedIn} from "../codegen/tables/UsedIn.sol";
import {UnitType} from "../codegen/tables/UnitType.sol";
import {AttackDamage} from "../codegen/tables/AttackDamage.sol";
import {ActionReady} from "../codegen/tables/ActionReady.sol";
import {MaxHp} from "../codegen/tables/MaxHp.sol";
import {CurrentHp} from "../codegen/tables/CurrentHp.sol";
import {MovementSpeed} from "../codegen/tables/MovementSpeed.sol";
import {Position} from "../codegen/tables/Position.sol";
import {IsBase} from "../codegen/tables/IsBase.sol";
import {LibDefaults} from "../libs/LibCards.sol";
import {PlacedCards} from "../codegen/tables/PlacedCards.sol";

contract JoinMatchSystem is System {
    function createCards(bytes32 player, uint32 playerIndex, bytes32 gameKey) private {
        // We are not reading the last CardTypes value because it is the base
        for (uint256 j = 0; j < uint256(type(CardTypes).max); j++) {
            bytes32 cardKey = bytes32(keccak256(abi.encodePacked(block.number, player, gasleft(), j)));
            Card.set(cardKey, true);
            OwnedBy.set(cardKey, player);
            UnitType.set(cardKey, CardTypes(j));
            UsedIn.set(cardKey, gameKey);
            AttackDamage.set(cardKey, LibDefaults.attack(CardTypes(j)));
            MaxHp.set(cardKey, LibDefaults.health(CardTypes(j)));
            CurrentHp.set(cardKey, LibDefaults.health(CardTypes(j)));
            MovementSpeed.set(cardKey, LibDefaults.movement(CardTypes(j)));
            ActionReady.set(cardKey, true);
        }

        bytes32 baseTR = bytes32(keccak256(abi.encodePacked(block.number, player, gasleft(), playerIndex + 1001)));
        bytes32 baseTL = bytes32(keccak256(abi.encodePacked(block.number, player, gasleft(), playerIndex + 1002)));
        bytes32 baseBR = bytes32(keccak256(abi.encodePacked(block.number, player, gasleft(), playerIndex + 1003)));
        bytes32 baseBL = bytes32(keccak256(abi.encodePacked(block.number, player, gasleft(), playerIndex + 1004)));
        if (playerIndex == 1) {
            Position.set(baseTR, true, gameKey, 4, 1);
            Position.set(baseTL, true, gameKey, 5, 1);
            Position.set(baseBR, true, gameKey, 4, 0);
            Position.set(baseBL, true, gameKey, 5, 0);
        } else {
            Position.set(baseTR, true, gameKey, 4, 9);
            Position.set(baseTL, true, gameKey, 5, 9);
            Position.set(baseBR, true, gameKey, 4, 8);
            Position.set(baseBL, true, gameKey, 5, 8);
        }

        bytes32 baseKey = bytes32(keccak256(abi.encodePacked(block.number, player, gasleft(), playerIndex + 1000)));
        IsBase.set(baseTR, baseKey);
        IsBase.set(baseTL, baseKey);
        IsBase.set(baseBR, baseKey);
        IsBase.set(baseBL, baseKey);
        // TODO: should we add usedIn for each part?

        // Set base values
        OwnedBy.set(baseKey, player);
        MaxHp.set(baseKey, 10);
        CurrentHp.set(baseKey, 10);
        UnitType.set(baseKey, CardTypes.Base);
        UsedIn.set(baseKey, gameKey);
    }

    function joinmatch(bytes32 key) public {
        bool value = Match.get(key);
        require(value, "match not found");

        bytes32 player1 = PlayerOne.get(key);
        require(PlayerOne.get(key) != 0, "player 1 is not set");

        bytes32 player2 = addressToEntityKey(_msgSender());
        require(PlayerTwo.get(key) == 0, "player 2 already set");
        // TODO: uncomment this so the player 2 is not the same as player 1
        require(player2 != player1, "player 1 and 2 must be different");
        PlayerTwo.set(key, player2);

        // Board config
        CurrentTurn.set(key, uint32(0));
        CurrentPlayer.set(key, player1);
        CurrentMana.set(key, uint32(5));
        PlacedCards.set(key, 0, 0);

        // Player Cards
        createCards(player1, 1, key);
        createCards(player2, 2, key);
    }
}
