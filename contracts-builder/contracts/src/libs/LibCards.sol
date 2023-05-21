// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import {CardTypes} from "../codegen/Types.sol";

library LibDefaults {
    function health(CardTypes cardType) internal pure returns (uint32) {
        if (cardType == CardTypes.VaanStrife) {
            return 6;
        }
        if (cardType == CardTypes.Felguard) {
            return 6;
        }
        if (cardType == CardTypes.Sakura) {
            return 5;
        }
        if (cardType == CardTypes.Freya) {
            return 5;
        }
        if (cardType == CardTypes.Lyra) {
            return 5;
        }
        if (cardType == CardTypes.Madmartigan) {
            return 10;
        }
        return 0;
    }

    function attack(CardTypes cardType) internal pure returns (uint32) {
        if (cardType == CardTypes.VaanStrife) {
            return 4;
        }
        if (cardType == CardTypes.Felguard) {
            return 4;
        }
        if (cardType == CardTypes.Sakura) {
            return 3;
        }
        if (cardType == CardTypes.Freya) {
            return 3;
        }
        if (cardType == CardTypes.Lyra) {
            return 1;
        }
        if (cardType == CardTypes.Madmartigan) {
            return 2;
        }
        return 0;
    }

    function movement(CardTypes cardType) internal pure returns (uint32) {
        if (cardType == CardTypes.VaanStrife) {
            return 2;
        }
        if (cardType == CardTypes.Felguard) {
            return 2;
        }
        if (cardType == CardTypes.Sakura) {
            return 3;
        }
        if (cardType == CardTypes.Freya) {
            return 3;
        }
        if (cardType == CardTypes.Lyra) {
            return 1;
        }
        if (cardType == CardTypes.Madmartigan) {
            return 1;
        }
        return 0;
    }
}
