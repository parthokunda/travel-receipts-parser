package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type receiptFile struct {
	fileName string
	path     string
	content  []byte
}

var parseCommand = &cobra.Command{
	Use:	"parse",
	Short: 	"parse receipts",
	Long: 	"parse your travel receipts using LLM",
	Run: 	func(cmd *cobra.Command, args []string) {
		parse("./Travel")
	},
}

func parse(srcPath string) {
	receipts, _ := get_files(srcPath)
	for _, receipt := range receipts {
		fmt.Println(receipt.path)
	}

	llmResponse, err := getDataFromReceipts(receipts)
	if err != nil {
		log.Fatal("LLM Parse Error")
	}

	fmt.Println(llmResponse)
}

func ExecuteCommands() {
	if err := parseCommand.Execute(); err != nil {
		log.Fatal("Failed to execute commands")
	}
}

func main() {
	ExecuteCommands()
	// receipts, _ := get_files("./Travel")
	// for _, receipt := range receipts {
	// 	fmt.Println(receipt.path)
	// }

	// llmResponse, err := getDataFromReceipts(receipts)
	// if err != nil {
	// 	log.Fatal("LLM Parse Error")
	// }

	// fmt.Println(llmResponse)

	// renameFilesUsingResponseAndMoveToProcessFolder(receipts, llmResponse)
	// srcPath := "Backup_Travel/receipt_8eba0e02-8e68-4b23-9871-a648cb7df4a4.pdf"
	// destPath := "Travel/test.pdf"
	// copyAndRenameFile(srcPath, destPath)

	// calculateTotalFareForMonth("07-2024")

}