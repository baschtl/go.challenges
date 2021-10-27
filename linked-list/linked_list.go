package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	Head *Node
}

func (l *LinkedList) Append(value int) {
	node := &Node{Value: value}

	if l.Head == nil {
		l.Head = node

		return
	}

	current := l.Head
	for current.Next != nil {
		current = current.Next
	}

	current.Next = node
}

func (l *LinkedList) Prepend(value int) {
	node := &Node{Value: value}

	if l.Head == nil {
		l.Head = node

		return
	}

	node.Next = l.Head
	l.Head = node
}

func (l *LinkedList) Delete(value int) {
	var previous *Node
	current := l.Head

	for current != nil {
		if current.Value == value {
			if previous == nil {
				l.Head = current.Next
				return
			} else {
				previous.Next = current.Next
				return
			}
		}
		previous = current
		current = current.Next
	}
}

func (l *LinkedList) Size() int {
	size := 0
	current := l.Head

	for current != nil {
		size++
		current = current.Next
	}

	return size
}

func (l *LinkedList) String() string {
	b := strings.Builder{}
	current := l.Head

	for current != nil {
		b.WriteString(fmt.Sprintf("%d ", current.Value))
		current = current.Next
	}

	return b.String()
}

func main() {
	l := LinkedList{}
	fmt.Println("Content: " + l.String())
	fmt.Println(fmt.Sprintf("Size: %d", l.Size()))

	fmt.Println("Appending some nodes...")

	l.Append(1)
	l.Append(3)
	l.Append(5)

	fmt.Println("Content: " + l.String())
	fmt.Println(fmt.Sprintf("Size: %d", l.Size()))

	fmt.Println("Prepending some nodes...")

	l.Prepend(-1)
	l.Prepend(-3)
	l.Prepend(-5)

	fmt.Println("Content: " + l.String())
	fmt.Println(fmt.Sprintf("Size: %d", l.Size()))

	fmt.Println("Deleting some nodes...")

	l.Delete(-5)
	l.Delete(-1)
	l.Delete(5)

	fmt.Println("Content: " + l.String())
	fmt.Println(fmt.Sprintf("Size: %d", l.Size()))
}
