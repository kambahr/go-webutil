package webutil

import (
	"fmt"
	"net/http"
	"strings"
)

// GenerateSessionID creates a sessionID that is comprised of a
// double MD5 of caller's IP address and user agent.
func (h *HTTP) GenerateSessionID(r *http.Request) string {

	sessionID := ""
	callerIPAddress := parseIPAddress(r)
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))

	sessionID = fmt.Sprintf("%s%s", callerIPAddress, userAgent)
	sessionID = getMd5String(getMd5String(sessionID))

	return sessionID
}

// GetCookie finds a target cookie and returns the string.
func (h *HTTP) GetCookie(cname string, r *http.Request) string {

	cookies := r.Cookies()

	for i := 0; i < len(cookies); i++ {

		if cookies[i].Name == cname {
			return cookies[i].Value
		}
	}

	return ""
}

// RemoveCookie a cookie by setting its expiration in the past.
func (h *HTTP) RemoveCookie(cname string, r *http.Request, w http.ResponseWriter) bool {

	cookies := r.Cookies()

	for i := 0; i < len(cookies); i++ {

		if cookies[i].Name != cname {
			continue
		}
		cx := fmt.Sprintf("%s=;path=/;expires=Thu, 01 Jan 1970 00:00:00;", cname)
		w.Header().Set("Set-Cookie", cx)
		return true
	}

	return false
}

// SetCookie creates a cookie by setting the Set-Cookie header. A value of zeor for
// maxAge will mean that the cookie will never expire.
func (h *HTTP) SetCookie(cname string, cvalue string, maxAge int, w http.ResponseWriter) string {

	maxAgeSec := 0

	if maxAge == 0 {
		maxAgeSec = (60 * 60 * 24 * 365 * 50) // Never expire
	} else {
		maxAgeSec = maxAge * 60
	}
	cx := fmt.Sprintf("%s=%s;path=/;max-age=%d;", cname, cvalue, maxAgeSec)
	w.Header().Set("Set-Cookie", cx)
	return cx
}
