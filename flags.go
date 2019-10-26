package main

import (
	//"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

// Usage: urlenc [flags] [input_file]
//       --bufsize int     I/O buffer size (default 1024)
//   -d, --decode          Decode input
//       --keepdelim       Don't remove delims
//       --ldelim string   Lefthand region delimiter (default ";;;")
//   -l, --line            Newline passthrough
//   -o, --output string   Output file (default: stdout)
//       --rdelim string   Righthand region delimiter (default ";;;")
//       --region          Encode or decode data within a
//                         predefined region that is defined
//                         by a lefthand (--ldelim) and
//                         righthand (--rdelim) boundary
//   -r, --rounds int      Encode/decode 'r' times (default 1)
//       --version         Print version
//   -h, --help            Print usage
//
//   input_file            optional (default: stdout)

var (
	fRounds     int
	fBufSize    int
	fLineMode   bool
	fKeepDelim  bool
	fInputPath  string
	fOutputPath string
	fLDelim     string
	fRDelim     string
	run         = encode
)

// init flags
func init() {
	var fDecodeMode, fRegionMode, fVersionMode bool

	// overwrite Usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "urlenc version %s\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [output_file]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  -h, --help            Print usage\n")
		fmt.Fprintf(os.Stderr, "\n  input_file            optional (default: stdout)\n")
	}

	// flags
	flag.BoolVarP(&fDecodeMode, "decode", "d", false, "Decode input")
	flag.BoolVarP(&fLineMode, "line", "l", false, "Newline passthrough")
	flag.BoolVar(&fRegionMode, "region", false, "Encode or decode data within a\n"+
		"predefined region that is defined\n"+
		"by a lefthand (--ldelim) and\n"+
		"righthand (--rdelim) boundary")
	flag.BoolVar(&fKeepDelim, "keepdelim", false, "Don't remove delims")
	flag.BoolVar(&fVersionMode, "version", false, "Print version")
	flag.IntVarP(&fRounds, "rounds", "r", 1, "Encode/decode 'r' times")
	flag.IntVar(&fBufSize, "bufsize", 1024, "I/O buffer size")
	flag.StringVarP(&fOutputPath, "output", "o", "", "Output file (default: stdout)")
	flag.StringVar(&fLDelim, "ldelim", ";;;", "Lefthand region delimiter")
	flag.StringVar(&fRDelim, "rdelim", ";;;", "Righthand region delimiter")

	// parse, run value sanity checks and set run function
	flag.Parse()
	if fVersionMode {
		fmt.Printf("urlenc %s\n", version)
		os.Exit(0)
	}
	flagSanity()
	setRunMode(fDecodeMode, fRegionMode)
}

// setRunMode determines the function to call by examining mode variables
func setRunMode(decodeMode, regionMode bool) {
	if !decodeMode && regionMode {
		run = encodeRegion
	} else if decodeMode && !regionMode {
		run = decode
	} else if decodeMode && regionMode {
		run = decodeRegion
	}
}

// flagSanity corrects wrong flag usage by resetting to default values
func flagSanity() {
	if fRounds < 1 {
		fRounds = 1
	}
	if fBufSize < 1 {
		fBufSize = 1024
	}
	if flag.NArg() >= 1 {
		fInputPath = flag.Arg(0)
	}
}
