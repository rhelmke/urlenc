# urlenc

[![Build Status](https://travis-ci.org/rhelmke/urlenc.svg?branch=master)](https://travis-ci.org/rhelmke/urlenc)
[![Go Report Card](https://goreportcard.com/badge/github.com/rhelmke/urlenc)](https://goreportcard.com/report/github.com/rhelmke/urlenc)

urlenc is a small and pretty fast commandline utility to URL-encode or -decode
selected datastreams.

While common tools like [urlencode(1)](https://linux.die.net/man/1/urlencode)
exist, urlenc provides extended functionality to double-, triple- or N-encode
streams. Furthermore, the `region` mode enables isolated encoding of predefined
regions. Both of these features may be of use when you want to process a large
list of URLs containing payloads.

urlenc does not adhere to [RFC3986](https://tools.ietf.org/html/rfc3986) and,
thus, encodes arbitrary bytes. Because why not?
[![utf8](https://i.imgur.com/XvXhSHL.gif)](https://asciinema.org/a/nJr0spUerGWErQfo1quPLOjvT?autoplay=1)

## Install

You can install urlenc via `go get`:

```bash
$ go get -u github.com/rhelmke/urlenc
```

But don't forget to put your go installation path (usually `$HOME/go/bin`) into
your `$PATH`!

Another option would be to download a prebuilt package from the
[Releases](https://github.com/rhelmke/urlenc/releases)-tab of this repository.

## Usage

```plain
Usage: urlenc [flags] [input_file]
      --bufsize int     I/O buffer size (default 1024)
  -d, --decode          Decode input
      --keepdelim       Don't remove delims
      --ldelim string   Lefthand region delimiter (default ";;;")
  -l, --line            Newline passthrough
  -o, --output string   Output file (default: stdout)
      --rdelim string   Righthand region delimiter (default ";;;")
      --region          Encode or decode data within a
                        predefined region that is defined
                        by a lefthand (--ldelim) and
                        righthand (--rdelim) boundary
  -r, --rounds int      Encode/decode 'r' times (default 1)
      --version         Print version
  -h, --help            Print usage

  input_file            optional (default: stdout)
```

In its default behavior, urlenc takes an input data stream from `stdin` and
writes each byte to `stdout`. If you want to alter this behavior, specify an
`input_file` or set an output file via the `-o` flag.

When setting the `-l` flag, urlenc does not encode newline characters:
[![line](https://i.imgur.com/c3LGmFl.gif)](https://asciinema.org/a/1rydKbsRjiv5O9bMKmx9t1yvv?autoplay=1&theme=monokai)
The `-r` flag can be used to specify the rounds a given byte will be encoded.
This can be used to, e.g., double-encode data:
[![dencode](https://i.imgur.com/C68YfrK.gif)](https://asciinema.org/a/APmSG87UrhI2RZRCbR5LxWAmu?autoplay=1&theme=monokai)
Of course, we can also decode data by using the `-d` flag:
[![decode](https://i.imgur.com/O57BdnX.gif)](https://asciinema.org/a/mFUd6lG9y3udzNDui8nKx4Tr2?autoplay=1&theme=monokai)
All flags that have been shown before do also apply to the decoding feature.

### Regions

Imagine you have a large list of URLs. Maybe there is a set of query parameters
you want to test for common vulnerabilities? But some of these characters need
encoding to be properly handled by the webserver. This is where regions come
into play.

By setting the `--region` flag, you can explicitly tell urlenc which parts of
a datastream should be encoded or decoded. The boundaries of a region are
defined by two strings that represent its beginning (`--ldelim`) and end
(`--rdelim`). You can define as many regions as you want within your input data
stream. Suppose your `payloads.txt` looks like this:

```plain
http://example.net/search?q=<script>alert(123);</script>
http://example.net/search?q=<ScRipT>alert("XSS");</ScRipT>
http://example.net/search?q=<script>alert(123)</script>
http://example.net/search?q=<script>alert("hellox worldss");</script>
http://example.net/search?q=<script>alert(“XSS”)</script>
http://example.net/search?q=<script>alert(“XSS”);</script>
http://example.net/search?q=<script>alert(‘XSS’)</script>
http://example.net/search?q=“><script>alert(“XSS”)</script>
http://example.net/search?q=<script>alert(/XSS”)</script>
http://example.net/search?q=<script>alert(/XSS/)</script>
```

Then you can insert a lefthand delimter, e.g. `^START^`, and a righthand delimiter,
e.g. `^STOP^` like this:

```plain
http://example.net/search?q=^START^<script>alert(123);</script>^STOP^
http://example.net/search?q=^START^<ScRipT>alert("XSS");</ScRipT>^STOP^
http://example.net/search?q=^START^<script>alert(123)</script>^STOP^
http://example.net/search?q=^START^<script>alert("hellox worldss");</script>^STOP^
http://example.net/search?q=^START^<script>alert(“XSS”)</script>^STOP^
http://example.net/search?q=^START^<script>alert(“XSS”);</script>^STOP^
http://example.net/search?q=^START^<script>alert(‘XSS’)</script>^STOP^
http://example.net/search?q=^START^“><script>alert(“XSS”)</script>^STOP^
http://example.net/search?q=^START^<script>alert(/XSS”)</script>^STOP^
http://example.net/search?q=^START^<script>alert(/XSS/)</script>^STOP^
```

Finally, pipe it to urlenc:
[![regions](https://i.imgur.com/YUCIciB.gif)](https://asciinema.org/a/JRHjIYiwhVIbkkn2LLqn4H8rc?autoplay=1&theme=monokai)

----

Oh, by the way. Since urlenc is entirely written in go, it supports UTF-8:
[![utf8](https://i.imgur.com/kodr9qG.gif)](https://asciinema.org/a/hd09IPZP6AV0XgAXvQHu8dYMz?autoplay=1&theme=monokai)
