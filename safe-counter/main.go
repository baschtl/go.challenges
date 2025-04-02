package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SafeCounterV1 struct {
	sync.Mutex
	value int
}

func (c *SafeCounterV1) Increment() {
	defer c.Unlock()

	c.Lock()
	c.value++
}

func (c *SafeCounterV1) Decrement() {
	defer c.Unlock()

	c.Lock()
	c.value--
}

func (c *SafeCounterV1) Get() int {
	return c.value
}

type SafeCounterV2 struct {
	value atomic.Int32
}

func (c *SafeCounterV2) Increment() {
	c.value.Add(1)
}

func (c *SafeCounterV2) Decrement() {
	c.value.Add(-1)
}

func (c *SafeCounterV2) Get() int {
	return int(c.value.Load())
}

func main() {
	fmt.Printf("Running Counter with mutex.\n")
	wg := sync.WaitGroup{}
	counterv1 := SafeCounterV1{}

	doIncrement := func() {
		counterv1.Increment()
		fmt.Printf("Counter: %d\n", counterv1.Get())
		wg.Done()
	}
	doDecrement := func() {
		counterv1.Decrement()
		fmt.Printf("Counter: %d\n", counterv1.Get())
		wg.Done()
	}

	wg.Add(6)
	for range 3 {
		go doIncrement()
		go doDecrement()
	}

	wg.Wait()
	fmt.Printf("Final Counter: %d\n\n", counterv1.Get())

	fmt.Printf("Running Counter with atomic type.\n")
	counterv2 := SafeCounterV2{}

	doIncrement = func() {
		counterv2.Increment()
		fmt.Printf("Counter: %d\n", counterv2.Get())
		wg.Done()
	}
	doDecrement = func() {
		counterv2.Decrement()
		fmt.Printf("Counter: %d\n", counterv2.Get())
		wg.Done()
	}

	wg.Add(6)
	for range 3 {
		go doIncrement()
		go doDecrement()
	}

	wg.Wait()
	fmt.Printf("Final Counter: %d\n", counterv2.Get())
}
