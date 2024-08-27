package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type receiptFile struct {
	fileName string
	path     string
	content  []byte
	fileType string
}


func get_files(dir string) ([]receiptFile, error) {
	var receipts []receiptFile
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Could not read from directory %s\n", dir)
			os.Exit(1)
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Could not read file %s", path)
				os.Exit(1)
			}
			// log.Println(filepath.Ext(path))
			receipts = append(receipts, receiptFile{info.Name(), path, content, filepath.Ext(path)})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func renameFilesUsingResponseAndCopyToProcessFolder(receipts []receiptFile, llmResponse string) {
	newNameForFiles := strings.Split(llmResponse, ";")

	processedFolderDir := "Travel/Processed/"
	for index, newFileName := range newNameForFiles {
		date := strings.Split(newFileName, "_")[0]
		folderName := processedFolderDir + strings.SplitN(date, "-", 2)[1] + "/"
		os.MkdirAll(folderName, os.ModePerm)
		
		destFilePath := folderName + newFileName + receipts[index].fileType
		srcFilePath := receipts[index].path

		copyAndRenameFile(srcFilePath, destFilePath)
	}
}

func copyAndRenameFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Printf("Failed to copy file %s to path %s", srcPath, destPath)
		return err
	}

	return nil
}