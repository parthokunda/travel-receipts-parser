package main


type receiptFile struct {
	fileName	string
	path 		string
	content		[]byte
}


func main() {
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
	srcPath := "Backup_Travel/receipt_8eba0e02-8e68-4b23-9871-a648cb7df4a4.pdf"
	destPath := "Travel/test.pdf"
	copyAndRenameFile(srcPath, destPath)


	// calculateTotalFareForMonth("07-2024")
	
}