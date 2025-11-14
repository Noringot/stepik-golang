package main

import "io"

func main() {
	for i := 0; i < 100000; i++ {
		FastSearch(io.Discard)
	}
	// fmt.Println()
	// fmt.Println("------------")
	// fmt.Println()
	// SlowSearch(os.Stdout)
}
