package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var directoryToParse string
var fileToParse string
var defaultDir string

var parseCommand = &cobra.Command{
	Use:   "parse",
	Short: "Parse Receipts",
	Long:  fmt.Sprintf("Parse your Travel Receipts (default \"%s\")", defaultDir),
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("file") {
			parse(fileToParse)
		} else if cmd.Flags().Changed("dir") {
			parse(directoryToParse)
		} else {
			parse(defaultDir)
		}
	},
}


func parse(srcPath string) {
	receipts, _ := get_files(srcPath)
	for _, receipt := range receipts {
		// fmt.Println(receipt.path)
		fmt.Println(receipt.fileType)
	}

	llmResponse, err := getDataFromReceipts(receipts)
	if err != nil {
		log.Fatal("LLM Parse Error")
	}

	fmt.Println(llmResponse)
	err = responseSanityCheck(llmResponse, receipts)
	if err != nil {
		log.Fatal(err)
	}
	renameFilesUsingResponseAndCopyToProcessFolder(receipts, llmResponse)
}

var calculateCommand = &cobra.Command{
	Use:	"calc",
	Short:	"Calculate your expenses",
	Run: func(cmd *cobra.Command, args []string) {
		calculateTotalFareForMonth("08-2024")
	},
}

func init() {
	defaultDir = "./Travel/Unprocessed"
	rootCmd.AddCommand(calculateCommand)

	parseCommand.Flags().StringVarP(&directoryToParse, "dir", "d", "./Travel/Unprocessed", "Directory to parse");
	parseCommand.Flags().StringVarP(&fileToParse, "file", "f", "", "File to parse")
	rootCmd.AddCommand(parseCommand)
}