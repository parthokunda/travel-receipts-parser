package main

import (
	"fmt"
	"os"
	"path/filepath"

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
			receipts = append(receipts, receiptFile{info.Name(), path, content})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	
	return receipts, nil
}

func main() {
	receipts, _ := get_files("./Travel")
	for _, receipt := range receipts {
		fmt.Println(receipt.path)
	}
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	ctx := context.Background()
	apiKey := os.Getenv("API_KEY")
	llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	parts := []llms.ContentPart{}
	
	for _, receipt := range receipts {
		parts = append(parts, llms.BinaryPart("application/pdf", receipt.content))
	}

	parts = append(parts, 
		llms.TextPart(
			`You have to extract the total fare(ignore currency) and the trip date from pdf. 
			Please output the information in the following format: DD-MM-YYYY_TotalFare. 
			For example if the fare is BDT 242.31 and Date 01/01/10 then out output should be 01-01-10_242.31.
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

	if len(resp.Choices) <= 0 {
		fmt.Println("No message received")
	} else {
		for _, choice := range resp.Choices {
			if len(choice.Content) > 0 {
				fmt.Println(choice.Content)
			}
		}
	}
	
}