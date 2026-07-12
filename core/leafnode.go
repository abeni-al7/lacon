package core

type LeafNode struct {
	letter string
	weight int
}

func NewLeafNode(letter string, weight int) LeafNode {
	return LeafNode{letter: letter, weight: weight}
}

func (l LeafNode) Letter() string {
	return l.letter
}

func (l LeafNode) Weight() int {
	return l.weight
}

func (l LeafNode) IsLeaf() bool {
	return true
}