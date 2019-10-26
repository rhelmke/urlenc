package main

import (
	"bufio"
	"io"
)

// decode all incoming data. Pass data through if not valid encoding
func decode(r *bufio.Reader, w *bufio.Writer) {
	var err error
	var b byte
	var n int
	defer w.Flush()      // explicit flush at end of stream
	t := []byte{0, 0, 0} // buffer used for left-shifting
	offs := 0            // offset to fill t
	currRound := 0       // current round specified by -r/--rounds
OUTER:
	for {
		// try to read 3-offs bytes at once
		read := 0
		for read != 3-offs {
			if n, err = r.Read(t[offs+read:]); err != nil {
				if err == io.EOF {
					// we may have to Write unhandled bytes after EOF
					if read > offs {
						offs = read
					}
					break OUTER
				}
				check(err, w)
			}
			read += n
		}
		// t probably contains a valid url code
		if t[0] == '%' {
			// try to decode
			if b, err = unhex(t[1], t[2]); err == nil {
				// decode ok!
				if b == '%' && currRound < fRounds-1 {
					//	in case we have multiple rounds to do, keep the decoded
					// '%' character in buffer
					t[0] = b
					offs = 1
					currRound++
					continue
				}
				// character is not '%' or we have no more rounds to do
				// pass character to output
				check(w.WriteByte(b), w)
				offs = 0
				currRound = 0
				continue
			}
		}
		// t probably doesn't contain a valid url code. Write first char in
		// t to output and shift t by one in order to get new characters from
		// stream
		check(w.WriteByte(t[0]), w)
		copy(t, t[1:])
		offs = 2
		currRound = 0
	}
	// flush everything that remains unhandled in t
	_, err = w.Write(t[:offs])
	check(err, w)
}
