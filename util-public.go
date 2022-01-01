package webutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// IsPortNoValid checks the ragne of an tcp/ip port number.
func (h *HTTP) IsPortNoValid(portno int) bool {
	return portno > 0 && portno <= 65535
}

// RemoveCommentsFromByBiteArry removes a block of text from a byte array.
func (h *HTTP) RemoveCommentsFromBiteArry(b []byte, begin string, end string) []byte {
	return h.RemoveCommentsFromByBiteArry(b, begin, end)
}

// RemoveCommentsFromByBiteArry removes a block of text from a byte array.
func (h *HTTP) RemoveCommentsFromByBiteArry(b []byte, begin string, end string) []byte {

	begin = strings.Trim(begin, " ")
	end = strings.Trim(end, " ")

	// Avoid recursion by using a label to go through
	// many iterations until all target blocks of text are removed.
lblAgain:
	i := bytes.Index(b, []byte(begin))

	if i < 0 {
		// not found
		return b
	}

	left := b[:i]

	right := b[len(left):]

	j := bytes.Index(right, []byte(end))
	right = right[j+len(end):]

	b = make([]byte, len(left)+len(right))

	k := 0
	for k = 0; k < len(left); k++ {
		b[k] = left[k]
	}

	for p := 0; p < len(right); p++ {
		b[k] = right[p]
		k++
	}

	i = bytes.Index(b, []byte(begin))

	if i > -1 {
		goto lblAgain
	}

	return b
}

//RemoveCommentsFromString removes a block of text from inside an string.
func (h *HTTP) RemoveCommentsFromString(s string, begin string, end string) string {

	begin = strings.Trim(begin, " ")
	end = strings.Trim(end, " ")

	// Avoid recursion by using a label to go through
	// many iterations until all target blocks of text are removed.
lblAgain:
	i := strings.Index(s, begin)

	if i < 0 {
		// not found
		return s
	}

	left := s[:i]

	right := s[len(left):]

	j := strings.Index(right, end)
	right = right[j+len(end):]

	s = fmt.Sprintf("%s%s", left, right)

	i = strings.Index(s, begin)

	if i > -1 {
		goto lblAgain
	}

	return s
}

// The following is the same as the Go ReadFile()
// func with the exception of closing the file before
// return.
//
// ../src/os/file.go
// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func ReadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		f.Close()
		return nil, err
	}

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			f.Close()
			return data, err
		}
	}
}
