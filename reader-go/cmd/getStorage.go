/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"reader-go/helpers"

	"github.com/aymantaybi/ronin/common"
	"github.com/aymantaybi/ronin/common/hexutil"
	"github.com/aymantaybi/ronin/crypto"
	"github.com/aymantaybi/ronin/trie"
	"github.com/spf13/cobra"
)

// getStorageCmd represents the getStorage command
var getStorageCmd = &cobra.Command{
	Use:   "getStorage",
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
		accountAddress := common.HexToAddress(args[1]) // ex: 0x0b7007c13325c48911f73a2dad5fa5dcbf808adc

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

		// The storage key (slot) you want to read
		storageKey := common.HexToHash(args[2]) // Provide the storage slot as the third argument ex: 0x5401e35c6516fe0403b686f8c632ea869cf7983f3f3e6eeea0fb274b009668e3

		// Ethereum storage trie uses the Keccak256 hash of the storage key
		hashedKey := crypto.Keccak256Hash(storageKey.Bytes())

		// Get the value from the storage trie
		storageValue, err := accountStorageTrie.TryGet(hashedKey.Bytes())
		if err != nil {
			fmt.Printf("Error getting storage value: %v\n", err)
			return
		}

		// Print the storage value
		fmt.Printf("Storage value at slot %s: %s\n", storageKey.Hex(), hexutil.Encode(storageValue))

	},
}

func init() {
	rootCmd.AddCommand(getStorageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getStorageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getStorageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
