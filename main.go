package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)


func main() {
	path := filepath.Join(".", "test", "test.txt")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	frequencyTable := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range line {
			frequencyTable[string(v)] += 1
		}
	}
	fmt.Print(frequencyTable)

}