package main

import "time"

type benchChanD struct {
	l    *bucket
	n    int
	a, b chan time.Time
}

func (bm *benchChanD) A() {
	for i := 0; i < bm.n; i++ {
		bm.a <- time.Now()
		s := <-bm.b
		bm.l.Record(time.Since(s))
	}
}

func (bm *benchChanD) B() {
	for i := 0; i < bm.n; i++ {
		s := <-bm.a
		bm.l.Record(time.Since(s))
		bm.b <- time.Now()
	}
}

func benchChan(l *bucket, n int) {
	b := &benchChanD{
		l,
		n,
		make(chan time.Time),
		make(chan time.Time),
	}

	go b.B()
	b.A()
}
