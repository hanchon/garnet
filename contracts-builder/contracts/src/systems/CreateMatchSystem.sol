// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import {System} from "@latticexyz/world/src/System.sol";
import {Match} from "../codegen/tables/Match.sol";
import {PlayerOne} from "../codegen/tables/PlayerOne.sol";
import {addressToEntityKey} from "../addressToEntityKey.sol";

contract CreateMatchSystem is System {
    function creatematch() public {
        bytes32 key = bytes32(abi.encodePacked(block.number, msg.sender, gasleft()));
        require(Match.get(key) != true, "game already exists");
        Match.set(key, true);

        bytes32 senderKey = addressToEntityKey(_msgSender());
        PlayerOne.set(key, senderKey);
    }
}
