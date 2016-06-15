package main

import (
	"encoding/binary"
	"io"
	"net"
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
	r.Close()
	w.Close()
}

func benchPipe(l *bucket, n int) {
	r, w, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	benchPipeGeneric(l, n, r, w)
	r.Close()
	w.Close()
}

func socketpipe() (net.Conn, net.Conn, error) {
	listen, err := net.ListenTCP("tcp6", &net.TCPAddr{net.IPv6loopback, 0, ""})
	if err != nil {
		return nil, nil, err
	}
	c := make(chan net.Conn)
	go func() {
		conn, err := listen.Accept()
		listen.Close()
		if err != nil {
			return
		}
		c <- conn
	}()
	r, err := net.Dial(listen.Addr().Network(), listen.Addr().String())
	if err != nil {
		listen.Close()
		return nil, nil, err
	}
	w := <-c
	close(c)

	r.(*net.TCPConn).SetNoDelay(true)
	w.(*net.TCPConn).SetNoDelay(true)

	return r, w, nil
}

func benchTcp(l *bucket, n int) {
	r, w, err := socketpipe()
	if err != nil {
		panic(err)
	}

	benchPipeGeneric(l, n, r, w)
	r.Close()
	w.Close()
}
