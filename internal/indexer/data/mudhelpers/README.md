# MudHelpers

To update the store contract, run on an already built mud lib:

```sh
cd mud/packages/store/abi/StoreCore.sol
abigen --abi storecore.abi.json --pkg mudhelpers --out generated.go
```
