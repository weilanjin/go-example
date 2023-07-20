package main

type Node[T any] struct {
	Value      T
	Prev, Next *Node[T]
}