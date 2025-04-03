package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

const numWorkers = 3

func worker(lines chan string, results chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	wordCount := make(map[string]int)

	for line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			trimmedWord := strings.ToLower(strings.Trim(word, ".,?!()'\""))
			wordCount[trimmedWord]++
		}
	}
	results <- wordCount
}

func merge(results chan map[string]int, final map[string]int, done chan struct{}) {
	for result := range results {
		for word, count := range result {
			final[word] += count
		}
	}
	done <- struct{}{}
}

func main() {
	file, err := os.Open("lorem.txt")
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup
	lines := make(chan string)           // Channel for sending lines to workers
	results := make(chan map[string]int) // Channel for sending work results of workers to result merge process
	mergeDone := make(chan struct{})     // Channel for indicating the end of the result merge process
	final := make(map[string]int)

	// Start worker goroutines
	for range numWorkers {
		wg.Add(1)
		go worker(lines, results, &wg)
	}

	// Start result merging goroutine
	go merge(results, final, mergeDone)

	// Read file line by line and send to workers
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	// Close the channel when all lines are sent
	close(lines)

	// Wait for workers to finish
	wg.Wait()
	// Close results after all workers are done
	close(results)

	// Wait for merge to complete
	<-mergeDone

	for word, count := range final {
		fmt.Printf("%s: %d\n", word, count)
	}
}
