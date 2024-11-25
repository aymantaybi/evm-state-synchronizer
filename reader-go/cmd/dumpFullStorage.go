/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"reader-go/helpers"

	"github.com/aymantaybi/ronin/common"
	"github.com/aymantaybi/ronin/common/hexutil"
	"github.com/aymantaybi/ronin/core/rawdb"
	"github.com/aymantaybi/ronin/core/types"
	"github.com/aymantaybi/ronin/crypto"
	"github.com/aymantaybi/ronin/rlp"
	"github.com/aymantaybi/ronin/trie"
	"github.com/spf13/cobra"
)

// dumpFullStorageCmd represents the dumpFullStorage command
var dumpFullStorageCmd = &cobra.Command{
	Use:   "dumpFullStorage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		stateRootInput := args[1]
		accountAddressInput := args[2]

		fmt.Printf("Access account (%v) state at root: %v\n", accountAddressInput, stateRootInput)

		handles := helpers.MakeDatabaseHandles(0)

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

		root, err := hexutil.Decode(stateRootInput) // "0xddc8f1b241f9090547501d92c2a943a41e8b076f14f2836be2cd8b4b1f6053c4"

		if err != nil {
			fmt.Printf("Error decoding root %v", err)
		}

		stateRoot := common.BytesToHash(root)

		theTrie, err := trie.New(stateRoot, trie.NewDatabase(db))

		if err != nil {
			fmt.Printf("Error creating storage state trie: %v\n", err)
		}

		// Address to retrieve
		addr := common.HexToAddress(accountAddressInput) // 0xc1eb47de5d549d45a871e32d9d082e7ac5d2e3ed
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
		fmt.Printf("Nonce: %d\n", account.Nonce)
		fmt.Printf("Balance: %s\n", account.Balance.String())
		fmt.Printf("Storage Root: %x\n", account.Root.Bytes())
		fmt.Printf("Code Hash: %x\n", account.CodeHash)

		// Create the storage trie using the account's storage root
		storageTrie, err := trie.New(account.Root, trie.NewDatabase(db))

		if err != nil {
			fmt.Printf("Error creating account storage trie: %v\n", err)
			return
		}

		// Iterate over the storage trie
		it := trie.NewIterator(storageTrie.NodeIterator(nil))
		var count int64
		for it.Next() {
			// Print the storage key and value
			count++
			fmt.Printf("Count: %d. key %#x: %#x\n", count, it.Key, it.Value)
		}

		fmt.Printf("Reading storage keys / values done!\n")

	},
}

func init() {
	rootCmd.AddCommand(dumpFullStorageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpFullStorageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpFullStorageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
