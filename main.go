package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "context"
	_ "log"

	_ "github.com/joho/godotenv"
	_ "github.com/tmc/langchaingo/llms"
	_ "github.com/tmc/langchaingo/llms/googleai"
)

type receiptFile struct {
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
			receipts = append(receipts, receiptFile{path, content})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	
	return receipts, nil
}

func main() {
	recepits, _ := get_files("./Travel")
	for _, receipt := range recepits {
		fmt.Println(receipt.path)
	}
	
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file:", err)
	// 	return
	// }

	// ctx := context.Background()
	// apiKey := os.Getenv("API_KEY")
	// llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pdfData, err := os.ReadFile("Travel/receipt_8eba0e02-8e68-4b23-9871-a648cb7df4a4.pdf")
	// if err != nil {
	// 	panic(err)
	// }

	// pdfData2, err := os.ReadFile("Travel/receipt_09f8563b-e242-4113-a239-717b9eb9dc87.pdf")
	// if err != nil {
	// 	panic(err)
	// }

	// parts := []llms.ContentPart{
	// 	llms.BinaryPart("application/pdf", pdfData),
	// 	llms.BinaryPart("application/pdf", pdfData2),
	// 	llms.TextPart(`I want to extract the total fare(ignore currency) and the trip date from pdf. Please output the information in the following format: TotalFare_DD-MM-YYYY. For example if the fare is BDT 242.31 and Date 01/01/10 then out output should be 242.31_01-01-10.`),
	// }

	// content := []llms.MessageContent{
	// 	{
	// 		Role:  llms.ChatMessageType("human"),
	// 		Parts: parts,
	// 	},
	// }

	// resp, err := llm.GenerateContent(ctx, content, llms.WithModel("gemini-1.5-flash-001"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // bs, _ := json.MarshalIndent(resp, "", "    ")
	// // fmt.Println(resp.Choices[0].Content)

	// if len(resp.Choices) <= 0 {
	// 	fmt.Println("No message received")
	// } else {
	// 	for _, choice := range resp.Choices {
	// 		if len(choice.Content) > 0 {
	// 			fmt.Println(choice.Content)
	// 		}
	// 	}
	// }
	
}