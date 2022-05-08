package webutil

import (
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
