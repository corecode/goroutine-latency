package main

import (
	"runtime"
	"time"
)

func benchEmpty(inf benchInfo) {
	c := make(chan struct{})
	go func() {
		if inf.locked {
			runtime.LockOSThread()
		}
		for i := 0; i < inf.n; i++ {
			t := time.Now()
			inf.l.Record(time.Since(t))
		}
		close(c)
	}()
	<-c
}
