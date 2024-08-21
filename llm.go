package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

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