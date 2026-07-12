package core

type Node interface {
	IsLeaf() bool
	Weight() int
}