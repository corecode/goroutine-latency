package main

import (
	"encoding/binary"
	"io"
	"os"
	"time"
)

func benchPipeGeneric(l *bucket, n int, r io.Reader, w io.Writer) {
	start := time.Now()

	go func() {
		for i := 0; i < n; i++ {
			time.Sleep(50 * time.Millisecond)
			d := time.Since(start)
			binary.Write(w, binary.LittleEndian, d)
		}
	}()

	for i := 0; i < n; i++ {
		var d time.Duration
		binary.Read(r, binary.LittleEndian, &d)
		l.Record(time.Since(start) - d)
	}
}

func benchPipeInternal(l *bucket, n int) {
	r, w := io.Pipe()

	benchPipeGeneric(l, n, r, w)
}

func benchPipe(l *bucket, n int) {
	r, w, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	benchPipeGeneric(l, n, r, w)
}
