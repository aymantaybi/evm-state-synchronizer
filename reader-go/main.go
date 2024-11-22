package main

import (
	"fmt"
	"path/filepath"

	"github.com/aymantaybi/ronin/common"
	"github.com/aymantaybi/ronin/common/fdlimit"
	"github.com/aymantaybi/ronin/common/hexutil"
	"github.com/aymantaybi/ronin/core/rawdb"
	"github.com/aymantaybi/ronin/core/types"
	"github.com/aymantaybi/ronin/crypto"
	"github.com/aymantaybi/ronin/rlp"
	"github.com/aymantaybi/ronin/trie"
)

func main() {

	handles := MakeDatabaseHandles(0)

	Directory := "../../roninchain/ronin/chaindata"

	freezer := filepath.Join(Directory, "ancient")

	db, err := rawdb.Open(rawdb.OpenOptions{
		Type:              "leveldb",
		Directory:         Directory,
		AncientsDirectory: freezer,
		Namespace:         "",
		Cache:             0,
		Handles:           handles,
		ReadOnly:          true,
	})

	if err != nil {
		fmt.Printf("Error opening database %v", err)
	}

	root, err := hexutil.Decode("0xddc8f1b241f9090547501d92c2a943a41e8b076f14f2836be2cd8b4b1f6053c4")

	if err != nil {
		fmt.Printf("Error decoding root %v", err)
	}

	stateRoot := common.BytesToHash(root)

	theTrie, err := trie.New(stateRoot, trie.NewDatabase(db))

	if err != nil {
		fmt.Printf("Error opening database %v", err)
	}

	// Address to retrieve
	addr := common.HexToAddress("0xc1eb47de5d549d45a871e32d9d082e7ac5d2e3ed")
	key := crypto.Keccak256(addr.Bytes())

	// Get the account data from the trie
	value := theTrie.Get(key)

	// Decode the RLP-encoded account data
	var account types.StateAccount
	err = rlp.DecodeBytes(value, &account)
	if err != nil {
		fmt.Printf("Error decoding account data: %v", err)
		return
	}

	// Print account details
	fmt.Printf("Account data for address %s:\n", addr.Hex())
	fmt.Printf("  Nonce: %d\n", account.Nonce)
	fmt.Printf("  Balance: %s\n", account.Balance.String())
	fmt.Printf("  Storage Root: %x\n", account.Root.Bytes())
	fmt.Printf("  Code Hash: %x\n", account.CodeHash)

	/* var count int64
	it := trie.NewIterator(theTrie.NodeIterator(nil))
	for it.Next() {
		fmt.Printf("  %d. key %#x: %#x\n", count, it.Key, it.Value)
		count++
	} */

}

// MakeDatabaseHandles raises out the number of allowed file handles per process
// for Geth and returns half of the allowance to assign to the database.
func MakeDatabaseHandles(max int) int {
	limit, err := fdlimit.Maximum()
	if err != nil {
		fmt.Println("Failed to retrieve file descriptor allowance: %v", err)
	}
	switch {
	case max == 0:
		// User didn't specify a meaningful value, use system limits
	case max < 128:
		// User specified something unhealthy, just use system defaults
		fmt.Println("File descriptor limit invalid (<128)", "had", max, "updated", limit)
	case max > limit:
		// User requested more than the OS allows, notify that we can't allocate it
		fmt.Println("Requested file descriptors denied by OS", "req", max, "limit", limit)
	default:
		// User limit is meaningful and within allowed range, use that
		limit = max
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		fmt.Println("Failed to raise file descriptor allowance: %v", err)
	}
	return int(raised / 2) // Leave half for networking and other stuff
}
