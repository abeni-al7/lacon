package main

type HuffmanTree struct {
	root Node
	weight int
}

func NewHuffmanTree(root *Node, weight int) HuffmanTree {
	return HuffmanTree{root: *root, weight: weight}
}

func (h HuffmanTree) Weight() int {
	return h.root.Weight()
}

func (h HuffmanTree) IsLeaf() bool {
	return h.root.IsLeaf()
}