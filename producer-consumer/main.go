package main

import (
	"fmt"
	"sync"
)

func main() {
	const jobs = 20
	const workers = 3

	ch := make(chan int)
	wg := sync.WaitGroup{}

	producer := func() {
		for i := range jobs {
			ch <- i
		}

		close(ch)
	}

	consumer := func(id int) {
		for i := range ch {
			fmt.Printf("Worker %d processed %d\n", id, i)
			wg.Done()
		}
	}

	wg.Add(jobs)
	go producer()

	for i := range workers {
		go consumer(i)
	}

	wg.Wait()
	fmt.Println("Done.")
}
