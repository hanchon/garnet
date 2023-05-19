import { mudConfig, resolveTableId } from "@latticexyz/world/register";

export default mudConfig({
  systems: {
    CreateMatchSystem: {
      name: "creatematch",
      openAccess: true,
    },
    JoinMatchSystem: {
      name: "joinmatch",
      openAccess: true,
    },
    PlaceCardSystem: {
      name: "placecard",
      openAccess: true,
    },
  },
  tables: {
    MapConfig: {
      primaryKeys: {},
      schema: {
        width: "uint32",
        height: "uint32",
        maxPlacedCards: "uint32",
      },
    },
    Match: "bool",
    PlayerOne: "bytes32",
    PlayerTwo: "bytes32",

    // Board
    CurrentTurn: "uint32",
    CurrentPlayer: "bytes32",
    CurrentMana: "uint32",
    PlacedCards: {
      schema: {
        p1Cards: "uint32",
        p2Cards: "uint32",
      },
    },

    // Units
    Card: "bool",
    OwnedBy: "bytes32",
    UsedIn: "bytes32", // relation between match and card
    IsBase: "bytes32",
    UnitType: "CardTypes",
    AttackDamage: "uint32",
    MaxHp: "uint32",
    CurrentHp: "uint32",
    MovementSpeed: "uint32",
    ActionReady: "bool",
    Position: {
      dataStruct: false,
      schema: {
        placed: "bool",
        gameKey: "bytes32",
        x: "uint32",
        y: "uint32",
      },
    },
  },

  enums: {
    // Base MUST be the last value
    CardTypes: ["Warrior1", "Warrior2", "Tank", "Mage", "Rogue1", "Rogue2", "Base"],
  },
  modules: [
    {
      name: "KeysWithValueModule",
      root: true,
      args: [resolveTableId("Position")],
    },
  ],
});
