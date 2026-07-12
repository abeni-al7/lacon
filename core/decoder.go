package core

import (
	"bufio"
	"encoding/json"
	"io"
)

func readAndRebuildPrefixTable(headerBytes []byte) (map[string]string, int, error) {
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

func decodeFromReader(input io.Reader, output io.Writer) error {
	reader := bufio.NewReader(input)

	headerBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	prefixTable, totalCount, err := readAndRebuildPrefixTable(headerBytes)
	if err != nil {
		return err
	}

	decodingMap := make(map[string]string)
	for char, bitString := range prefixTable {
		decodingMap[bitString] = char
	}

	bitReader := NewBitReader(reader)
	outputWriter := bufio.NewWriter(output)

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

// Decode decompresses data from r and writes the decompressed output to w.
func Decode(r io.Reader, w io.Writer) error {
	return decodeFromReader(r, w)
}