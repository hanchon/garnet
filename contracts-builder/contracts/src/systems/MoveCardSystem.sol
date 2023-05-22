// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// Core
import {System} from "@latticexyz/world/src/System.sol";
import {getKeysWithValue} from "@latticexyz/world/src/modules/keyswithvalue/getKeysWithValue.sol";
// Utils
import {addressToEntityKey} from "../addressToEntityKey.sol";
import {CardTypes} from "../codegen/Types.sol";
// Tables
import {Card} from "../codegen/tables/Card.sol";
import {OwnedBy} from "../codegen/tables/OwnedBy.sol";
import {UsedIn} from "../codegen/tables/UsedIn.sol";
import {CurrentPlayer} from "../codegen/tables/CurrentPlayer.sol";
import {PlayerOne} from "../codegen/tables/PlayerOne.sol";
import {UnitType} from "../codegen/tables/UnitType.sol";
import {Position, PositionTableId} from "../codegen/tables/Position.sol";
import {PlacedCards, PlacedCardsData} from "../codegen/tables/PlacedCards.sol";
import {MapConfig, MapConfigData} from "../codegen/tables/MapConfig.sol";
import {CurrentMana} from "../codegen/tables/CurrentMana.sol";
import {MovementSpeed} from "../codegen/tables/MovementSpeed.sol";

contract MoveCardSystem is System {
    function validate(bytes32 cardKey, bytes32 gameKey, bytes32 playerKey) private view {
        require(Card.get(cardKey), "card does not exist");
        require(OwnedBy.get(cardKey) == playerKey, "the sender is not the owner of the card");
        require(CurrentPlayer.get(gameKey) == playerKey, "it is not the player turn");
    }

    function checkEmptyPos(bytes32 gameKey, uint32 x, uint32 y) private view {
        // Check that there is no card in that position
        bytes32[] memory keysAtPos = getKeysWithValue(PositionTableId, Position.encode(true, gameKey, x, y));
        require(keysAtPos.length == 0, "there is a unit in that position");
    }

    function limits(bytes32 cardKey, bytes32 gameKeyGenerated, bytes32 playerKey, uint32 newX, uint32 newY)
        private
        view
        returns (PlacedCardsData memory)
    {
        // Is the card the base
        require(UnitType.get(cardKey) != CardTypes.Base, "can not place move the base");

        // Is the card in the board?
        (bool placed,, uint32 x, uint32 y) = Position.get(cardKey);
        require(placed == true, "card was not summoned");

        // Check max distance using movement speed
        uint32 deltaX = newX > x ? newX - x : x - newX;
        uint32 deltaY = newY > y ? newY - y : y - newY;
        require(deltaX + deltaY <= MovementSpeed.get(cardKey), "the card is trying to move too far");

        // Map limits
        MapConfigData memory mapConfig = MapConfig.get();
        require(newX <= mapConfig.width && newX >= 0, "invalid x");
        require(newY <= mapConfig.height && newY >= 0, "invalid y");

        // Mana
        require(CurrentMana.get(gameKeyGenerated) >= 2, "no enough mana");
    }

    function movecard(bytes32 cardKey, uint32 newX, uint32 newY) public {
        bytes32 gameKeyGenerated = UsedIn.get(cardKey);
        bytes32 playerKey = addressToEntityKey(_msgSender());
        require(gameKeyGenerated != 0, "game id is incorrect");
        validate(cardKey, gameKeyGenerated, playerKey);

        // Check that there is no card in that position
        checkEmptyPos(gameKeyGenerated, newX, newY);
        // limits
        limits(cardKey, gameKeyGenerated, playerKey, newX, newY);

        // Update game status
        Position.set(cardKey, true, gameKeyGenerated, newX, newY);
        CurrentMana.set(gameKeyGenerated, CurrentMana.get(gameKeyGenerated) - 2);
    }
}
