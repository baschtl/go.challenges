package main

import (
	"fmt"
	"strings"
)

type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

type LRUCache struct {
	lookupMap map[int]*Node
	head      *Node
	tail      *Node
	capacity  int
	size      int
}

func NewLRUCache(capacity int) *LRUCache {
	lru := &LRUCache{
		lookupMap: make(map[int]*Node),
		capacity:  capacity,
		size:      0,
		head:      &Node{},
		tail:      &Node{},
	}

	lru.head.next = lru.tail
	lru.tail.prev = lru.head

	return lru
}

func (lru *LRUCache) Print() string {
	builder := strings.Builder{}

	currentNode := lru.head
	for currentNode != nil {
		builder.WriteString(fmt.Sprintf("%d: %d", currentNode.key, currentNode.value))

		if currentNode.next != nil {
			builder.WriteString(" -> ")
		}
		currentNode = currentNode.next
	}

	return builder.String()
}

func (lru *LRUCache) Get(key int) int {
	if node, found := lru.lookupMap[key]; found {
		lru.moveToFront(node)

		return node.value
	} else {
		return -1
	}
}

func (lru *LRUCache) Put(key, value int) {
	// Handle update case
	if node, found := lru.lookupMap[key]; found {
		node.value = value
		lru.moveToFront(node)

		return
	}

	// If capacity is reached we need to remove the last node before adding a new one
	if lru.size == lru.capacity {
		lru.popTail()
	}

	// Handle add new node case
	node := &Node{
		key:   key,
		value: value,
	}

	lru.lookupMap[key] = node
	lru.moveToFront(node)
	lru.size++
}

func (lru *LRUCache) moveToFront(node *Node) {
	// Remove node from its original spot
	lru.removeNode(node)

	// Move node after head node
	node.next = lru.head.next
	node.prev = lru.head
	lru.head.next.prev = node
	lru.head.next = node
}

func (lru *LRUCache) popTail() *Node {
	toPop := lru.tail.prev

	lru.removeNode(toPop)
	delete(lru.lookupMap, toPop.key)
	lru.size--

	return toPop
}

func (lru *LRUCache) removeNode(node *Node) {
	if node.prev == nil || node.next == nil {
		return
	}

	node.prev.next = node.next
	node.next.prev = node.prev
}

func main() {
	lru := NewLRUCache(2)

	lru.Put(1, 1)
	lru.Put(2, 2)
	fmt.Println(lru.Print())

	lru.Put(3, 3)
	fmt.Println(lru.Print())

	lru.Put(4, 4)
	fmt.Println(lru.Print())

	lru.Put(3, 3)
	fmt.Println(lru.Print())

	fmt.Printf("Found node %d\n", lru.Get(4))
	fmt.Println(lru.Print())
}
