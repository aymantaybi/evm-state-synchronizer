/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"reader-go/helpers"

	"github.com/aymantaybi/ronin/common"
	"github.com/aymantaybi/ronin/common/hexutil"
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

		directory := "../../roninchain/ronin/chaindata"

		db, err := helpers.OpenRawDB(directory)
		if err != nil {
			fmt.Printf("Error opening database %v\n", err)
			return
		}

		root, err := hexutil.Decode(args[0]) // ex: 0xddc8f1b241f9090547501d92c2a943a41e8b076f14f2836be2cd8b4b1f6053c4
		if err != nil {
			fmt.Printf("Error decoding root %v\n", err)
			return
		}

		stateTrie, err := helpers.GetStateTrie(db, root)
		if err != nil {
			fmt.Printf("Error getting storage state trie: %v\n", err)
		}

		// Address to retrieve
		accountAddress := common.HexToAddress(args[1]) // ex: 0xc1eb47de5d549d45a871e32d9d082e7ac5d2e3ed

		accountState, err := helpers.GetAccountState(stateTrie, accountAddress)
		if err != nil {
			fmt.Printf("Error getting account state: %v", err)
			return
		}

		// Print account details
		fmt.Printf("Account state for address %s:\n", accountAddress.Hex())
		fmt.Printf("Nonce: %d\n", accountState.Nonce)
		fmt.Printf("Balance: %s\n", accountState.Balance.String())
		fmt.Printf("Storage Root: %x\n", accountState.Root.Bytes())
		fmt.Printf("Code Hash: %x\n", accountState.CodeHash)

		// Create the storage trie using the account's storage root
		accountStorageTrie, err := trie.New(accountState.Root, trie.NewDatabase(db))
		if err != nil {
			fmt.Printf("Error creating account storage trie: %v\n", err)
			return
		}

		// Iterate over the storage trie
		it := trie.NewIterator(accountStorageTrie.NodeIterator(nil))
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
