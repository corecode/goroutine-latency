package main

import "time"

func benchEmpty(l *bucket, n int) {
	for i := 0; i < n; i++ {
		t := time.Now()
		l.Record(time.Since(t))
	}
}
