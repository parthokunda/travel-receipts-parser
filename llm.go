package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

func getBinaryPart(receipt receiptFile) llms.BinaryContent {
	if (receipt.fileType == "jpg") {
		return llms.BinaryPart("application/jpg", receipt.content)
	} else if(receipt.fileType == "png") {
		return llms.BinaryPart("application/png", receipt.content)
	} else if(receipt.fileType == "pdf") {
		return llms.BinaryPart("application/pdf", receipt.content)
	} else {
		return llms.BinaryPart(http.DetectContentType(receipt.content), receipt.content)
	}
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
		parts = append(parts, llms.BinaryPart(http.DetectContentType(receipt.content), receipt.content))
		fmt.Println(http.DetectContentType(receipt.content))
	}

	parts = append(parts,
		llms.TextPart(
			`I am going to upload some travel receipts from ride sharing platforms.
			You have to extract the cost I have incurred for the ride and the trip date from all files. 
			Please output the information in the following format: DD-MM-YYYY_TotalFare. 
			For example if the fare is BDT 242.31 and Date 01/01/10 then its output should be 01-01-10_242.31.
			For multiple files, separate each one with semicolon(;). 
			Each entry should be in sequence as I have uploaded them to you.
			If you fail to parse any file, only output "I am noob".
			For two files as input the output could be 01-01-10_242.31;31-01-11_192.14`))

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

func responseSanityCheck(resp string, receipts []receiptFile) error {
	newNameForFiles := strings.Split(resp, ";")
	if(len(newNameForFiles) != len(receipts)){
		return errors.New("failed to get response for all files.")
	}
	pattern := `^(\d{2}-\d{2}-\d{4}_\d+(\.\d+)?)(;\d{2}-\d{2}-\d{4}_\d+(\.\d+)?)*$`
	regex := regexp.MustCompile(pattern)

	if ! regex.MatchString(strings.TrimSpace(resp)) {
		return errors.New("invalid response format from llm")
	}
	return nil
}