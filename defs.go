// Package webutil implements some Web utility functions.
// It also implments a cache mechanism for any file that is
// to be served via http.
package webutil

import (
	"time"

	"github.com/kambahr/go-webcache"
)

const cachetypes = ".js .css .min.css .min.css.map .js.map .min.js .min.js.map .csv .xls .xlsx .ods"

// HTTP are common http callback functions.
type HTTP struct {
	RootPhysicalPath string
	CacheDuration    time.Duration
	Webcache         *webcache.Cache
}

// UserSession defines a web user session.
type UserSession struct {
	Email                        string
	CustomerID                   int
	FirstName                    string
	LastName                     string
	EmailVerified                bool
	SessionID                    string
	VerificationCodeExpires      time.Time
	LastVerificationCodeRequest  time.Time
	VerificationCodeRequestCount int
}

//go:generate stringer -type=HTTPMethod
type HTTPMethod int

const (
	GET HTTPMethod = iota
	HEAD
	POST
	PUT
	CONNECT
	DELETE
	OPTIONS
	PATCH
	TRACE
)
