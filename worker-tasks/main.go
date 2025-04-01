package main

import (
	"fmt"
	"sync"
)

// Write a Go program that:
//     Spawns N worker goroutines.
//     Receives M tasks (numbers from 1 to M).
//     Each worker squares the number and prints the result.
//     Use channels to distribute tasks and synchronize completion.

// For example, with N=3 workers and M=5 tasks, output could be:

// Worker 2 processing 2 → 4
// Worker 1 processing 1 → 1
// Worker 3 processing 3 → 9
// Worker 2 processing 4 → 16
// Worker 1 processing 5 → 25
// All tasks completed.

func worker(id int, ch chan (int), wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d waiting...\n", id)
	for t := range ch {
		fmt.Printf("Worker %d processing %d → %d\n", id, t, t*t)
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	workers := 3
	tasks := 5

	for i := range workers {
		wg.Add(1)
		go worker(i+1, ch, &wg)
	}

	for i := range tasks {
		ch <- i + 1
	}
	close(ch)

	wg.Wait()
	fmt.Println("All tasks completed.")
}
