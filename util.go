package webutil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// parseIPAddress gets the left side of the r.RemoteAddr.
// It also checks for the local IP6 addr.
func parseIPAddress(r *http.Request) string {

	ipAddress := r.RemoteAddr
	if strings.Contains(ipAddress, "::") {
		ipAddress = "127.0.0.1"
	} else {
		values := strings.Split(ipAddress, ":")
		ipAddress = values[0]
	}

	return ipAddress
}

// fileOrDirectoryExists checks if a file or directory exists.
func fileOrDirectoryExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// getMd5String computes the MD5 and returns string format.
func getMd5String(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// getFileExtension gets full extension of a path.
func getFileExtension(p string) string {

	extx := ""

	if p == "" {
		return ""
	}

	v := strings.Split(p, "/")

	if len(v) == 0 {
		return ""
	}

	t := v[len(v)-1]

	d := strings.Split(t, ".")

	if len(d) < 2 {
		return ""
	}

	for i := 1; i < len(d); i++ {
		extx = fmt.Sprintf("%s.%s", extx, d[i])
	}

	return extx
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
