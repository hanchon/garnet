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

contract PlaceCardSystem is System {
    function validateCard(bytes32 cardKey, bytes32 gameKey, bytes32 playerKey) private view {
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
        // Just allow placing cards in the first 2 rows
        require(
            (playerKey == PlayerOne.get(gameKeyGenerated) && (newY == 0 || newY == 1))
                || (playerKey != PlayerOne.get(gameKeyGenerated) && (newY == 8 || newY == 9)),
            "incorrect row"
        );

        // Is the card the base
        require(UnitType.get(cardKey) != CardTypes.Base, "can not place the base");

        // Is the card in the board?
        (bool placed,,,) = Position.get(cardKey);
        require(placed == false, "card already placed");

        // Map limits
        MapConfigData memory mapConfig = MapConfig.get();
        require(newX <= mapConfig.width && newX >= 0, "invalid x");
        require(newY <= mapConfig.height && newY >= 0, "invalid y");

        PlacedCardsData memory placedCards = PlacedCards.get(gameKeyGenerated);
        // No more than 3 cards can be played
        if (playerKey == PlayerOne.get(gameKeyGenerated)) {
            require(placedCards.p1Cards < mapConfig.maxPlacedCards, "already placed the max amount of cards");
            placedCards.p1Cards = placedCards.p1Cards + 1;
        } else {
            require(placedCards.p2Cards < mapConfig.maxPlacedCards, "already placed the max amount of cards");
            placedCards.p2Cards = placedCards.p2Cards + 1;
        }

        // Mana
        require(CurrentMana.get(gameKeyGenerated) >= 3, "no enough mana");
        return placedCards;
    }

    function placecard(bytes32 cardKey, uint32 newX, uint32 newY) public {
        bytes32 gameKeyGenerated = UsedIn.get(cardKey);
        bytes32 playerKey = addressToEntityKey(_msgSender());
        require(gameKeyGenerated != 0, "game id is incorrect");
        validateCard(cardKey, gameKeyGenerated, playerKey);

        // Check that there is no card in that position
        checkEmptyPos(gameKeyGenerated, newX, newY);
        // limits
        PlacedCardsData memory placedCards = limits(cardKey, gameKeyGenerated, playerKey, newX, newY);

        // Update game status
        Position.set(cardKey, true, gameKeyGenerated, newX, newY);
        PlacedCards.set(gameKeyGenerated, placedCards.p1Cards, placedCards.p2Cards);
        CurrentMana.set(gameKeyGenerated, CurrentMana.get(gameKeyGenerated) - 3);
    }
}
