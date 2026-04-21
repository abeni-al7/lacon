package main

type Node interface {
	IsLeaf() bool
	Weight() int
}