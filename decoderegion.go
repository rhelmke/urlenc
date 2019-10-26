package main

import (
	"bufio"
	"io"
)

// decodeRegion only decodes data that is inside a predefined region.
// e.g.: javascript:;;;%61%6C%65%72%74%28%31%29;;; => javascript:alert(1)
func decodeRegion(r *bufio.Reader, w *bufio.Writer) {
	var err error
	var b byte
	var n int
	var delimSize int
	var delim string
	delimCnt := 0
	// we need a history for bytes we already read because we won't
	// necessarily write them immediately to the output
	delimHistory := make([]byte, maxInt(len(fLDelim), len(fRDelim)))
	insideRegion := false
	defer w.Flush()      // explicit flush at end of stream
	t := []byte{0, 0, 0} // buffer used for left-shifting
	offs := 0            // offset to fill t
	currRound := 0       // current round specified by -r/--rounds
OUTER:
	for {
		// check each iteration if inside region and set variables as
		// needed
		if insideRegion {
			delimSize = len(fRDelim)
			delim = fRDelim
		} else {
			delimSize = len(fLDelim)
			delim = fLDelim
		}
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
		for i := 0; i < read; i++ {
			if t[i] == delim[delimCnt] {
				// this byte is one of the bytes we are looking for to get
				// the defined delimiter.
				delimHistory[delimCnt] = b
				delimCnt++
				if delimCnt == delimSize {
					// we are now switching modes since history == delimiter
					if fKeepDelim {
						_, err = w.Write(delimHistory[:delimCnt])
						check(err, w)
					}
					insideRegion = !insideRegion
					delimCnt = 0
				}
				continue
			}
			delimCnt = 0
		}
		// t probably contains a valid url code
		if t[0] == '%' && insideRegion {
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
		// t to output and shift t by one in order to get new characters in
		// the next iteration step
		check(w.WriteByte(t[0]), w)
		copy(t, t[1:])
		offs = 2
		currRound = 0
	}
	// flush everything that remains unhandled in t
	_, err = w.Write(t[:offs])
	check(err, w)
}
