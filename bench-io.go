package main

import (
	"encoding/binary"
	"io"
	"net"
	"os"
	"runtime"
	"time"
)

func benchPipeGeneric(inf benchInfo, r io.Reader, w io.Writer) {
	start := time.Now()

	go func() {
		if inf.locked {
			runtime.LockOSThread()
		}
		for i := 0; i < inf.n; i++ {
			time.Sleep(50 * time.Millisecond)
			d := time.Since(start)
			binary.Write(w, binary.LittleEndian, d)
		}
	}()

	c := make(chan struct{})
	go func() {
		if inf.locked {
			runtime.LockOSThread()
		}
		for i := 0; i < inf.n; i++ {
			var d time.Duration
			binary.Read(r, binary.LittleEndian, &d)
			inf.l.Record(time.Since(start) - d)
		}
		close(c)
	}()
	<-c
}

func benchPipeInternal(inf benchInfo) {
	r, w := io.Pipe()

	benchPipeGeneric(inf, r, w)
	r.Close()
	w.Close()
}

func benchPipe(inf benchInfo) {
	r, w, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	benchPipeGeneric(inf, r, w)
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

func benchTcp(inf benchInfo) {
	r, w, err := socketpipe()
	if err != nil {
		panic(err)
	}

	benchPipeGeneric(inf, r, w)
	r.Close()
	w.Close()
}
