package helpers

import (
	"fmt"

	"github.com/aymantaybi/ronin/common"
	"github.com/aymantaybi/ronin/core/types"
	"github.com/aymantaybi/ronin/crypto"
	"github.com/aymantaybi/ronin/ethdb"
	"github.com/aymantaybi/ronin/rlp"
	"github.com/aymantaybi/ronin/trie"
)

func GetStateTrie(db ethdb.Database, root []byte) (*trie.Trie, error) {
	stateRoot := common.BytesToHash(root)
	return trie.New(stateRoot, trie.NewDatabase(db))
}

func GetAccountState(stateTrie *trie.Trie, address common.Address) (types.StateAccount, error) {
	key := crypto.Keccak256(address.Bytes())
	value := stateTrie.Get(key)
	var account types.StateAccount
	err := rlp.DecodeBytes(value, &account)
	if err != nil {
		fmt.Printf("Error decoding account data: %v", err)
	}
	return account, err
}
