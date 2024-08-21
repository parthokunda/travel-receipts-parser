package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getFareFromFileName(filename string) float64 {
	filenameWithoutExt := strings.Join(strings.Split(filename, ".")[0:2], ".") // there are two dot, so needed bit of hack
	fareString := strings.Split(filenameWithoutExt, "_")[1]
	fare, err := strconv.ParseFloat(fareString, 64)
	if err != nil {
		log.Fatal(err)
	}

	return fare
}

func calculateTotalFareForMonth(month string) {
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
