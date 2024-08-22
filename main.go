package main

import (
	"log"

	"github.com/spf13/cobra"
)

type receiptFile struct {
	fileName string
	path     string
	content  []byte
}

var rootCmd = &cobra.Command{
	Use:   "tr",
	Short: "Extract information from travel receipts",
	Long:  "A program to organize and extract information from your travel receipts",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func ExecuteCommands() {
	if err := rootCmd.Execute(); err != nil {
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