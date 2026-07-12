package core

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"io"
)

func countCharacterOccurrences(reader io.Reader) (map[string]int, error) {
	frequencyTable := make(map[string]int)

	bufReader := bufio.NewReader(reader)
	for {
		char, _, err := bufReader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		frequencyTable[string(char)] += 1
	}

	return frequencyTable, nil
}

func buildHuffmanTree(frequencyTable map[string]int) HuffmanTree {
	h := &MinHeap{}
	heap.Init(h)

	for letter, freq := range frequencyTable {
		node := NewLeafNode(letter, freq)
		current := NewHuffmanTree(node)

		heap.Push(h, current)
	}

	for h.Len() > 1 {
		leftTree := heap.Pop(h).(HuffmanTree)
		rightTree := heap.Pop(h).(HuffmanTree)

		var left Node = leftTree.Root()
		var right Node = rightTree.Root()
		newNode := NewInternalNode(left.Weight()+right.Weight(), &left, &right)
		newTree := NewHuffmanTree(newNode)
		heap.Push(h, newTree)
	}

	return h.Pop().(HuffmanTree)
}

func constructPrefixCodeTable(table map[string]string, node HuffmanTree, current string) {
	if s, ok := node.Root().(LeafNode); ok {
		key := s.Letter()
		if current == "" {
			current = "0" // handle trees with single node
		}
		table[key] = current
	} else if s, ok := node.Root().(InternalNode); ok {
		constructPrefixCodeTable(table, NewHuffmanTree(s.Left()), current+"0")
		constructPrefixCodeTable(table, NewHuffmanTree(s.Right()), current+"1")
	}
}

func writeHeader(freqTable map[string]int, output io.Writer) error {
	jsonData, err := json.Marshal(freqTable)
	if err != nil {
		return err
	}

	_, err = output.Write(jsonData)
	if err != nil {
		return err
	}

	_, err = output.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}

func writeContents(prefixCodeTable map[string]string, input io.Reader, output io.Writer) error {
	bitWriter := NewBitWriter(output)

	reader := bufio.NewReader(input)
	for {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := bitWriter.Writebits(prefixCodeTable[string(char)]); err != nil {
			return err
		}
	}
	if err := bitWriter.Flush(); err != nil {
		return err
	}
	return nil
}

// Encode compresses data from r and writes the compressed output to w.
// r must be seekable so that two passes can be made over the data.
func Encode(r io.ReadSeeker, w io.Writer) error {
	// First pass: count character frequencies
	frequencyTable, err := countCharacterOccurrences(r)
	if err != nil {
		return err
	}

	// Build Huffman tree and prefix code table
	tree := buildHuffmanTree(frequencyTable)

	prefixCodeTable := make(map[string]string)
	constructPrefixCodeTable(prefixCodeTable, tree, "")

	// Write header (frequency table as JSON)
	if err := writeHeader(frequencyTable, w); err != nil {
		return err
	}

	// Rewind to the beginning for the second pass
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Second pass: write encoded contents
	if err := writeContents(prefixCodeTable, r, w); err != nil {
		return err
	}

	return nil
}
