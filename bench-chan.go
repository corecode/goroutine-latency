package main

import (
	"runtime"
	"time"
)

type benchChanD struct {
	benchInfo
	a, b chan time.Time
}

func (bm *benchChanD) A() {
	if bm.locked {
		runtime.LockOSThread()
	}
	for i := 0; i < bm.n; i++ {
		bm.a <- time.Now()
		s := <-bm.b
		bm.l.Record(time.Since(s))
	}
}

func (bm *benchChanD) B() {
	if bm.locked {
		runtime.LockOSThread()
	}
	for i := 0; i < bm.n; i++ {
		s := <-bm.a
		bm.l.Record(time.Since(s))
		bm.b <- time.Now()
	}
}

func benchChan(inf benchInfo) {
	b := &benchChanD{
		benchInfo: inf,
		a:         make(chan time.Time),
		b:         make(chan time.Time),
	}

	c := make(chan struct{})
	go b.B()
	go func() {
		b.A()
		close(c)
	}()
	<-c
}
