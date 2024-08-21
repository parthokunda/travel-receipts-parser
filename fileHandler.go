package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
		return err
	}

	return nil
}