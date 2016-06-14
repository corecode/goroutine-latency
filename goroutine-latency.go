package main

import (
	"fmt"
	"time"
)

func main() {
	n := 0
	for {
		fmt.Printf("activity=%d\n", n)
		bench(benchEmpty)
		bench(benchChan)
		bench(benchPipe)
		bench(benchPipeInternal)

		for i := 0; i < 10; i++ {
			n++
			go activity()
		}
		time.Sleep(time.Second)
	}
}
