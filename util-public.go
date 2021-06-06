package webutil

import (
	"bytes"
	"fmt"
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

	left := b[:i-1]
	right := b[len(left):]

	j := bytes.Index(right, []byte(end))
	right = right[j+2:]

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

	left := s[:i-1]

	right := s[len(left):]

	j := strings.Index(right, end)
	right = right[j+2:]

	s = fmt.Sprintf("%s%s", left, right)

	i = strings.Index(s, begin)

	if i > -1 {
		goto lblAgain
	}

	return s
}
