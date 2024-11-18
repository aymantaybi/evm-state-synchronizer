# evm-state-synchronizer

EVM contracts state synchronizer

Starting from a geth level db node snapshot at certain height, load contracts of interest in an EVM (REVM), using the logs of the following blocks update the state of the contracts, save the latest state in JSON files.

merkle-patricia-tree: https://www.npmjs.com/package/merkle-patricia-tree (Other alternatives can be considered for performance, something in rust or go)

Using the merkle-patricia-tree package read the contract storage keys & values, save the pairs as JSON file.

Read the json file, load the contracts storage in an in-memory revm db, consume the logs of the next blocks until the latest block, using custom handlers to sync the in-memory db.
