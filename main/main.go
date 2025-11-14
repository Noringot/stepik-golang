package main

import (
	"fmt"
	"math/rand"
	"time"
)

func processData(val int) int {
	time.Sleep(time.Duration(rand.Intn(1)) * time.Second)
	return val * 2
}

func worker(in <-chan int, out chan<- int) {
	for v := range in {
		out <- processData(v)
	}
}

func processParallel(in <-chan int, out chan<- int, numWorkers int) {
	for range numWorkers {
		go worker(in, out)
	}

	select {
	case <-time.After(5 * time.Second):
		close(out)
		return
	}
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 100 {
			in <- i
		}

		close(in)
	}()

	now := time.Now()
	go processParallel(in, out, 5)

	for v := range out {
		fmt.Println(v)
	}

	fmt.Println(time.Since(now))
}
