package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// charset
var charset = "0123456789ABCDEF"

// will be thrown
var errUnhex = errors.New("can not unhex")

// maxInt determines the maximum of two given integers a and b since the std
// math package does not include an integer version
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// check if function returned an error. Clean up in case of error and and exit.
// this function exists because deferred statements are not executed when invoking
// os.Exit(..). Also, it reduces code size
func check(err error, w *bufio.Writer) {
	if err == nil {
		return
	}
	if w != nil {
		w.Flush()
	}
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

// hex receives an arbitrary byte and converts it to a (high, low)-printable
// hexadecimal representation
func hex(b byte) (byte, byte) {
	return charset[b>>0x04], charset[b&0x0F] // high and low nibbles
}

// unhex receives the bytes of a hex (high, low)-printable hexadecimal
// representation and converts them back to the original byte. This function
// returns errUnhex if a byte in either
func unhex(high, low byte) (byte, error) {
	var err error
	var hNorm, lNorm byte
	// convert the most significant nibble
	if hNorm, err = normalize(high); err != nil {
		return 0, err
	}
	// convert the least significant nibble
	if lNorm, err = normalize(low); err != nil {
		return 0, err
	}
	// join nibbles to a proper byte
	return (hNorm << 0x04) | lNorm, nil
}

// normalize a given byte that represents [0-9a-zA-Z] in ascii to its original
// value. This function is used in conjunction with unhex and returns errUnhex
// if b is not within the ascii range specified above.
func normalize(b byte) (byte, error) {
	switch {
	case '0' <= b && b <= '9':
		return b - '0', nil
	case 'a' <= b && b <= 'f':
		return b - 'a' + 10, nil
	case 'A' <= b && b <= 'F':
		return b - 'A' + 10, nil
	}
	// b is not within [0-9a-zA-Z] and can thus not be normalized
	return 0, errUnhex
}
