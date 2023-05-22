// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// Core
import {System} from "@latticexyz/world/src/System.sol";
import {getKeysWithValue} from "@latticexyz/world/src/modules/keyswithvalue/getKeysWithValue.sol";
// Tables
import {Match} from "../codegen/tables/Match.sol";
import {PlayerOne} from "../codegen/tables/PlayerOne.sol";
import {PlayerTwo} from "../codegen/tables/PlayerTwo.sol";
import {CurrentTurn} from "../codegen/tables/CurrentTurn.sol";
import {CurrentPlayer} from "../codegen/tables/CurrentPlayer.sol";
import {CurrentMana} from "../codegen/tables/CurrentMana.sol";
import {addressToEntityKey} from "../addressToEntityKey.sol";
import {UsedIn, UsedInTableId} from "../codegen/tables/UsedIn.sol";
import {ActionReady} from "../codegen/tables/ActionReady.sol";

contract EndTurnSystem is System {
    function updateCards(bytes32 matchKey) public {
        // Check that there is no card in that position
        bytes32[] memory cards = getKeysWithValue(UsedInTableId, UsedIn.encode(matchKey));
        require(cards.length != 0, "there are no units to update");
        for (uint256 j = 0; j < cards.length; j++) {
            // TODO: if type is base do not set this flag
            ActionReady.set(cards[j], true);
        }
    }

    function endturn(bytes32 key) public {
        bool value = Match.get(key);
        require(value, "match not found");

        bytes32 currentPlayer = CurrentPlayer.get(key);
        require(addressToEntityKey(_msgSender()) == currentPlayer, "current player must be the sender");

        bytes32 playerTwo = PlayerTwo.get(key);
        bytes32 playerOne = PlayerOne.get(key);
        uint32 ct = CurrentTurn.get(key);
        if (playerTwo == currentPlayer) {
            CurrentPlayer.set(key, playerOne);
        } else {
            CurrentPlayer.set(key, playerTwo);
        }

        CurrentTurn.set(key, ct + 1);

        // Mana
        if (ct + 5 + 1 > 15) {
            CurrentMana.set(key, uint32(15));
        } else {
            CurrentMana.set(key, uint32(ct + 5 + 1));
        }

        updateCards(key);
    }
}
