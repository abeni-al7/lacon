package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	inputPath := filepath.Join(".", "test", "test.txt")
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputPath := filepath.Join(".", "test", "test_output.txt")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	huffmanEncode(inputFile, outputFile)

	if _, err := outputFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	decoded, err := os.Create(filepath.Join(".", "test", "decoded.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer decoded.Close()

	huffmanDecode(outputFile, decoded)
}