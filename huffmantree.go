package main

type HuffmanTree struct {
	root Node
	weight int
}

func NewHuffmanTree(root *Node) HuffmanTree {
	return HuffmanTree{root: *root, weight: (*root).Weight()}
}

func (h HuffmanTree) Weight() int {
	return h.root.Weight()
}

func (h HuffmanTree) IsLeaf() bool {
	return h.root.IsLeaf()
}

func (h HuffmanTree) Less(other HuffmanTree) bool {
	return h.weight <= other.weight
}