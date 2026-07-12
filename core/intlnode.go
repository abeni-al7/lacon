package core

type InternalNode struct {
	weight int
	left   Node
	right  Node
}

func NewInternalNode(weight int, l *Node, r *Node) InternalNode {
	return InternalNode{weight: weight, left: *l, right: *r}
}

func (i InternalNode) Weight() int {
	return i.weight
}

func (i InternalNode) IsLeaf() bool {
	return false
}

func (i InternalNode) Left() Node {
	return i.left
}

func (i InternalNode) Right() Node {
	return i.right
}