package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type benchInfo struct {
	l      *bucket
	n      int
	locked bool
}

func bench(b func(inf benchInfo)) {
	inf := benchInfo{
		l:      nil,
		n:      100,
		locked: false,
	}

	for {
		fmt.Printf("%v...\n", runtime.FuncForPC(reflect.ValueOf(b).Pointer()).Name())
		for {
			start := time.Now()
			inf.l = newBucket()
			b(inf)

			duration := time.Since(start)

			if duration > 2*time.Second {
				break
			}

			if duration*10 > 2*time.Second {
				inf.n = int(int64(2*time.Second) * int64(inf.n) / int64(duration))
			} else {
				inf.n *= 10
			}
		}
		fmt.Printf("n = %d, locked = %v\n", inf.n, inf.locked)
		fmt.Printf("%v\n", inf.l)
		if inf.locked {
			break
		}
		inf.locked = true
	}
}
