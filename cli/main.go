package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/abeni-al7/lacon/core"
)

func main() {
	encode := flag.Bool("e", false, "encode the input file")
	decode := flag.Bool("d", false, "decode the input file")
	flag.Parse()

	if (*encode && *decode) || (!*encode && !*decode) {
		fmt.Fprintln(os.Stderr, "Error: specify exactly one of -e or -d")
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: input and output files are required")
		flag.Usage()
		os.Exit(1)
	}

	input := args[0]
	output := args[1]

	inputFile, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	if *encode {
		if err := core.Encode(inputFile, outputFile); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := core.Decode(inputFile, outputFile); err != nil {
			log.Fatal(err)
		}
	}
}