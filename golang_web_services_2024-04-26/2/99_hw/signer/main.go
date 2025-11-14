package main

import "fmt"

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})

	go func(in, out chan interface{}) {
		for v := range out {
			in <- v
		}
	}(in, out)

	for _, job := range jobs {
		go job(in, out)
	}

	select {}
}

func main() {
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			fmt.Println("Worker #1 write")
		}),
		job(func(in, out chan interface{}) {
			for v := range in {
				fmt.Printf("Worker #2 read value %v \n", v)
				out <- 12
				fmt.Println("Worker #2 write")
			}
		}),
		job(func(in, out chan interface{}) {
			for v := range in {
				fmt.Printf("Worker #3 read value %v \n", v)
			}
		}),
	}
	ExecutePipeline(freeFlowJobs...)
}
