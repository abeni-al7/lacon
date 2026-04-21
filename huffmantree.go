package main

type HuffmanTree struct {
	root Node
	left Node
	right Node
	weight int
}

func NewHuffmanTree(root *Node, l *Node, r *Node, weight int) HuffmanTree {
	return HuffmanTree{root: *root, left: *l, right: *r, weight: weight}
}

func (h HuffmanTree) Weight() int {
	return h.root.Weight()
}

func (h HuffmanTree) IsLeaf() bool {
	return h.root.IsLeaf()
}