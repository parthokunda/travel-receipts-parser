package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var parseCommand = &cobra.Command{
	Use:   "parse",
	Short: "parse receipts",
	Long:  "parse your travel receipts using LLM",
	Run: func(cmd *cobra.Command, args []string) {
		parse("./Backup_Travel")
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


func init() {
	rootCmd.AddCommand(parseCommand)
}