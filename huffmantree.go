package main

type HuffmanTree struct {
	root    Node
	weight  int
	sortKey string
}

func NewHuffmanTree(root Node) HuffmanTree {
	return HuffmanTree{root: root, weight: root.Weight(), sortKey: treeSortKey(root)}
}

func (h HuffmanTree) Weight() int {
	return h.root.Weight()
}

func (h HuffmanTree) IsLeaf() bool {
	return h.root.IsLeaf()
}

func (h HuffmanTree) Less(other HuffmanTree) bool {
	if h.weight != other.weight {
		return h.weight < other.weight
	}

	return h.sortKey < other.sortKey
}