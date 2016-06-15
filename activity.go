package main

import "crypto/sha256"

func activity() {
	r, w, _ := socketpipe() //os.Pipe()

	go func() {
		data := make([]byte, sha256.Size)
		for {
			r.Read(data)
			sha256.Sum256(data)
		}
	}()

	data := sha256.Sum256([]byte("hi"))
	for {
		data = sha256.Sum256(data[:])
		w.Write(data[:])
	}
}
