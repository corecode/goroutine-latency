package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func bench(b func(l *bucket, n int)) {
	n := 1000
	var l *bucket

	fmt.Printf("%v...\n", runtime.FuncForPC(reflect.ValueOf(b).Pointer()).Name())
	for {
		start := time.Now()
		l = newBucket()
		b(l, n)

		duration := time.Since(start)

		if duration > 2*time.Second {
			break
		}

		if duration*10 > 2*time.Second {
			n = int(int64(2*time.Second) * int64(n) / int64(duration))
		} else {
			n *= 10
		}
	}
	fmt.Printf("n = %d\n", n)
	fmt.Printf("%v\n", l)
}
