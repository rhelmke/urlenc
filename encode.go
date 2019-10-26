package main

import (
	"bufio"
	"io"
)

// encode all incoming data in url format (does not adhere to RFC 3986)
func encode(r *bufio.Reader, w *bufio.Writer) {
	var err error
	var b byte
	// explicit flush at end of stream
	defer w.Flush()
	// read byte-wise until error or EOF. bufio handles performance
	for {
		b, err = r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			check(err, w)
		}
		_enc(b, w)
	}
}

// _enc is an internal function that wraps the actual encoding of a byte
// and directly writes it to the bufio.Writer
func _enc(b byte, w *bufio.Writer) {
	var err error
	// this holds
	t := []byte{0, 0}
	// newline and lineMode -> passthrough
	if b == '\n' && fLineMode {
		check(w.WriteByte(b), w)
		return
	}
	// write '%' rune and optionally apply multiple rounds for n-encoding
	check(w.WriteByte('%'), w)
	t[0], t[1] = '2', '5'
	for i := 1; i < fRounds; i++ {
		_, err = w.Write(t)
		check(err, w)
	}
	// finally write the originally encoded input value
	t[0], t[1] = hex(b)
	_, err = w.Write(t)
	check(err, w)
}
