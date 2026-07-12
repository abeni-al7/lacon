package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

func ReadAndRebuildPrefixTable(headerBytes []byte) (map[string]string, int, error) {
	var frequencyTable map[string]int
	err := json.Unmarshal(headerBytes, &frequencyTable)
	if err != nil {
		return nil, 0, err
	}

	totalCount := 0
	for _, frequency := range frequencyTable {
		totalCount += frequency
	}

	tree := buildHuffmanTree(frequencyTable)
	prefixCodeTable := make(map[string]string)
	constructPrefixCodeTable(prefixCodeTable, tree, "")

	return prefixCodeTable, totalCount, nil
}

func huffmanDecode(inputFile *os.File, outputFile *os.File) error {
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	headerBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	prefixTable, totalCount, err := ReadAndRebuildPrefixTable(headerBytes)
	if err != nil {
		return err
	}

	decodingMap := make(map[string]string)
	for char, bitString := range prefixTable {
		decodingMap[bitString] = char
	}

	bitReader := NewBitReader(reader)
	outputWriter := bufio.NewWriter(outputFile)

	var currentBits []byte
	decodedCount := 0
	for {
		if decodedCount >= totalCount {
			break
		}

		bit, err := bitReader.ReadBit()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		currentBits = append(currentBits, byte(bit))

		if originalChar, exists := decodingMap[string(currentBits)]; exists {
			_, err := outputWriter.WriteString(originalChar)
			if err != nil {
				return err
			}
			decodedCount++
			currentBits = currentBits[:0]
		}
	}
	if err := outputWriter.Flush(); err != nil {
		return err
	}
	return nil
}