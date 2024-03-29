// Package webutil implements some Web utility functions.
// It also implments a cache mechanism for any file that is
// to be served via http.
package webutil

import (
	"net/http"
	"time"

	"github.com/kambahr/go-webcache"
)

const (
	PageCmdLoadFile = "$LoadFile"
)

const cachetypes = ".js .css .min.css .min.css.map .js.map .min.js .min.js.map .csv .xls .xlsx .ods"

// HTTP are common http callback functions.
type HTTP struct {
	RootPhysicalPath string
	CacheDuration    time.Duration
	Webcache         *webcache.Cache

	DecryptFile         func(filePath string, keyPhrase string)
	EncryptionKeyPhrase string

	// SecurityToken is passed via querystring to restrict
	// asset (static) file. If not blank every call must have
	// ?key=<secrity token> in order to access the file.
	SecurityToken string
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
	ReferredURL                  string
}

type HTTPResult struct {
	ResponseData []byte
	RequestDump  string
	Response     *http.Response
	Request      *http.Request
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
