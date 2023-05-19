// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import {CardTypes} from "../codegen/Types.sol";

library LibDefaults {
    function health(CardTypes cardType) internal pure returns (uint32) {
        if (cardType == CardTypes.Warrior1) {
            return 12;
        }
        if (cardType == CardTypes.Warrior2) {
            return 10;
        }
        if (cardType == CardTypes.Tank) {
            return 7;
        }
        if (cardType == CardTypes.Mage) {
            return 8;
        }
        if (cardType == CardTypes.Rogue1) {
            return 16;
        }
        if (cardType == CardTypes.Rogue2) {
            return 11;
        }
        return 0;
    }

    function attack(CardTypes cardType) internal pure returns (uint32) {
        if (cardType == CardTypes.Warrior1) {
            return 12;
        }
        if (cardType == CardTypes.Warrior2) {
            return 10;
        }
        if (cardType == CardTypes.Tank) {
            return 7;
        }
        if (cardType == CardTypes.Mage) {
            return 8;
        }
        if (cardType == CardTypes.Rogue1) {
            return 16;
        }
        if (cardType == CardTypes.Rogue2) {
            return 11;
        }
        return 0;
    }

    function movement(CardTypes cardType) internal pure returns (uint32) {
        if (cardType == CardTypes.Warrior1) {
            return 12;
        }
        if (cardType == CardTypes.Warrior2) {
            return 10;
        }
        if (cardType == CardTypes.Tank) {
            return 7;
        }
        if (cardType == CardTypes.Mage) {
            return 8;
        }
        if (cardType == CardTypes.Rogue1) {
            return 16;
        }
        if (cardType == CardTypes.Rogue2) {
            return 11;
        }
        return 0;
    }
}
