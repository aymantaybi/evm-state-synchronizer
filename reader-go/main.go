package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func main() {
	// Database path
	dbPath := "../../roninchain/ronin/chaindata" // Adjust this path

	// Open the LevelDB database
	db, err := leveldb.New(dbPath, 0, 0, "")
	if err != nil {
		log.Fatalf("Failed to open LevelDB database: %v", err)
	}
	defer db.Close()
	fmt.Println("Database is open and ready for operations.")

	// Create a trie database
	trieDb := trie.New(db)

	// State root hash
	stateRootHex := "0xddc8f1b241f9090547501d92c2a943a41e8b076f14f2836be2cd8b4b1f6053c4"
	stateRoot := common.HexToHash(stateRootHex)

	// Initialize the state trie
	trie, err := trie.New(stateRoot, trieDb)
	if err != nil {
		log.Fatalf("Failed to create state trie: %v", err)
	}

	// Account address
	addressHex := "0x0b7007c13325c48911f73a2dad5fa5dcbf808adc"
	address := common.HexToAddress(addressHex)

	// Get the account data
	addrHash := crypto.Keccak256(address.Bytes())
	accountRLP, err := trie.Get(addrHash)
	if err != nil {
		log.Fatalf("Failed to get account data: %v", err)
	}
	if accountRLP == nil {
		fmt.Println("Account not found in state trie.")
		return
	}

	// Decode the account
	var account state.Account
	if err := rlp.DecodeBytes(accountRLP, &account); err != nil {
		log.Fatalf("Failed to decode account data: %v", err)
	}

	// Display account information
	fmt.Println("-------State-------")
	fmt.Printf("Nonce: %d\n", account.Nonce)
	fmt.Printf("Balance in wei: %s\n", account.Balance.String())
	fmt.Printf("Storage Root: %s\n", account.Root.Hex())
	fmt.Printf("Code Hash: %s\n", common.BytesToHash(account.CodeHash).Hex())

	// Access the storage trie
	storageTrie, err := trie.New(account.Root, trieDb)
	if err != nil {
		log.Fatalf("Failed to create storage trie: %v", err)
	}

	// Iterate over storage trie
	fmt.Println("------Storage------")
	it := trie.NewIterator(storageTrie.NodeIterator(nil))
	for it.Next() {
		// Decode storage key and value
		key := it.Key
		value := it.Value

		// The storage key is hashed; to get the original key, reverse the hashing
		// which may not be feasible without knowing the original key.
		// We'll display the hashed key.

		fmt.Printf("Key Hash: %s\n", hex.EncodeToString(key))

		// Decode the value
		var storageValue big.Int
		if err := rlp.DecodeBytes(value, &storageValue); err != nil {
			log.Printf("Failed to decode storage value: %v", err)
			continue
		}
		fmt.Printf("Value: %s\n", storageValue.String())
	}
	fmt.Println("Finished reading storage.")
}
