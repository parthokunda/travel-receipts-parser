package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type receiptFile struct {
	fileName	string
	path 		string
	content		[]byte
}

func get_files(dir string) ([]receiptFile, error) {
	var receipts []receiptFile
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			fmt.Println(filepath.Ext(path))	
			receipts = append(receipts, receiptFile{info.Name(), path, content})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	
	return receipts, nil
}

func getDataFromReceipts(receipts []receiptFile) (string, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return "", err
	}

	ctx := context.Background()
	apiKey := os.Getenv("API_KEY")
	llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	parts := []llms.ContentPart{}
	
	for _, receipt := range receipts {
		parts = append(parts, llms.BinaryPart("application/pdf", receipt.content))
	}

	parts = append(parts, 
		llms.TextPart(
			`You have to extract the total fare(ignore currency) and the trip date from pdf. 
			Please output the information in the following format: DD-MM-YYYY_TotalFare. 
			For example if the fare is BDT 242.31 and Date 01/01/10 then its output should be 01-01-10_242.31.
			For multiple files, separate each one with semicolon(;). Keep the original file sequence.
			Example could be 01-01-10_242.31;31-01-11_192.14`))

	content := []llms.MessageContent{
		{
			Role:  llms.ChatMessageType("human"),
			Parts: parts,
		},
	}

	resp, err := llm.GenerateContent(ctx, content, llms.WithModel("gemini-1.5-flash-001"))
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Choices) != 1 {
		log.Fatal("response length should be one. but has length " + string(len(resp.Choices)))
	}

	return resp.Choices[0].Content, nil
}

func renameFilesUsingResponseAndMoveToProcessFolder(receipts []receiptFile, llmResponse string) {
	newNameForFiles := strings.Split(llmResponse, ";")
	processedFolderDir := "Travel/Processed/"
	for index, newFileName := range newNameForFiles {
		date := strings.Split(newFileName, "_")[0]
		folderName := processedFolderDir + strings.SplitN(date, "-", 2)[1] + "/"
		os.MkdirAll(folderName, os.ModePerm)
		
		err := os.Rename(receipts[index].path, folderName + newFileName + ".pdf")
		if err != nil {
			panic(err)
		}
	}
}

func getFareFromFileName(filename string) float64 {
	filenameWithoutExt := strings.Join(strings.Split(filename, ".")[0:2], ".") // there are two dot, so needed bit of hack
	fareString := strings.Split(filenameWithoutExt, "_")[1]
	fare, err := strconv.ParseFloat(fareString, 64)
	if err != nil {
		log.Fatal(err)
	}

	return fare
}

func calculateTotalFareForMonth(month string){
	monthFolderDir := "Travel/Processed/" + month
	directoryEntries, err := os.ReadDir(monthFolderDir)
	if err != nil {
		log.Fatal("Directory Not Found")
	}

	totalFare := 0.0
	for _, file := range directoryEntries {
		if !file.IsDir() {
			totalFare += (getFareFromFileName(file.Name()))
		}
	}
	fmt.Printf("Total Fare: %.2f\n", totalFare)

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


	calculateTotalFareForMonth("07-2024")
	
}