package main

import (
	"bufio"
	"io"
)

// encodeRegion only encodes data that is inside a predefined region.
// e.g.: javascript:;;;alert(1);;; => javascript:%61%6C%65%72%74%28%31%29
func encodeRegion(r *bufio.Reader, w *bufio.Writer) {
	var err error
	var b byte
	var delimSize int
	var delim string
	delimCnt := 0
	// we need a history for bytes we already read because we won't
	// necessarily write them immediately to the output
	delimHistory := make([]byte, maxInt(len(fLDelim), len(fRDelim)))
	insideRegion := false
	defer w.Flush()
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
		// read byte
		b, err = r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			check(err, w)
		}
		if b == delim[delimCnt] {
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
		if insideRegion {
			// only encode data within region
			for i := 0; i < delimCnt; i++ {
				_enc(delimHistory[i], w)
			}
			_enc(b, w)
		} else {
			// otherwise passthrough
			_, err = w.Write(delimHistory[:delimCnt])
			check(err, w)
			check(w.WriteByte(b), w)
		}
		delimCnt = 0
	}
}
