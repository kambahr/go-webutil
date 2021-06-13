package webutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/kambahr/go-webcache"
)

// NewHTTP creates a new instance of webutil.
func NewHTTP(rootPhysicalPath string, cacheDuration time.Duration) *HTTP {
	var h HTTP
	h.RootPhysicalPath = rootPhysicalPath
	h.CacheDuration = cacheDuration
	h.Webcache = webcache.NewWebCache(cacheDuration)
	return &h
}

// setContentTypeAndWrite writes the response and reutrns false, if mime type not found;
// returns true if mime type found.
// The returns are for info -- as in any case the mime type is written:
//    1. Mime type not found, let the browser handle it.
//    2. Mime type found but not chaced, write from the file.
//    3. Mime type found and cached, write from the cache.
func (h *HTTP) setContentTypeAndWrite(w http.ResponseWriter, r *http.Request) (bool, bool, bool) {

	mimTypeFound := false
	servedFromCache := false
	servedFromFile := false
	uriPath := r.URL.Path
	ext := getFileExtension(uriPath)

	cntType := mime.TypeByExtension(ext)
	if cntType != "" && !strings.Contains(cachetypes, ext) {
		// Let the browser/server handle the ones not on the list of cachetypes.
		w.Header().Set("Content-Type", cntType)
		return true, servedFromCache, servedFromFile
	}

	// All else fall into the cache category.
	cntType = h.GetMIMEContentType(ext)
	if cntType != "" {
		mimTypeFound = true
	}
	w.Header().Set("Content-Type", cntType)
	var b []byte
	var err error
	physPath := fmt.Sprintf("%s%s", h.RootPhysicalPath, uriPath)
	bFromCache := h.Webcache.GetItem(uriPath)
	if len(bFromCache) == 0 {
		b, err = ioutil.ReadFile(physPath)
		if err != nil {
			// TO DO: log this or notify the caller. 
			fmt.Println(err)
		} else {
			h.Webcache.AddItem(uriPath, b, h.CacheDuration)
			writeResponse(b, w , r)
			servedFromFile = true
		}
	} else {
		writeResponse(bFromCache, w , r)
		servedFromCache = true
	}

	return mimTypeFound, servedFromCache, servedFromFile
}

// SetContentTypeAndWrite writes the response and reutrns false, if mime type not found;
// returns true if mime type found. It uses the conent passed via an arg rather than
// than that of the request.
func (h *HTTP) SetContentTypeAndWrite(w http.ResponseWriter, r *http.Request, f []byte) bool {
	mimTypeFound := false
	uriPath := strings.ToLower(r.URL.Path)
	ext := getFileExtension(uriPath)

	cntType := mime.TypeByExtension(ext)
	if cntType != "" && !strings.Contains(cachetypes, ext) {
		// Let the browser/server handle the ones not on the list of cachetypes.
		w.Header().Set("Content-Type", cntType)
		return true
	}

	// All else fall into the cache category.
	cntType = h.GetMIMEContentType(ext)
	if cntType != "" {
		mimTypeFound = true
	}
	w.Header().Set("Content-Type", cntType)
	writeResponse(f, w , r)

	return mimTypeFound
}

// writeResponse compresses the response if the client asks for it.
// TODO: add more types of compressions.
func writeResponse(content []byte, w http.ResponseWriter, r *http.Request) {

	if r.Method == "HEAD" {
		return
	}

	compressResponse := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	if compressResponse {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Transfer-Encoding", "gzip, chunked")
	}

	w.WriteHeader(http.StatusOK)

	if compressResponse {
		var b bytes.Buffer
		gw, _ := gzip.NewWriterLevel(&b, gzip.DefaultCompression)
		gw.Write(content)
		gw.Close()
		w.Write(b.Bytes())
	} else {
		w.Write(content)
	}
}

// GetMIMEContentType first checks the standard extensions i.e. .png, .js,...
// if not found it uses a custom parsing to return the right content type.
func (h *HTTP) GetMIMEContentType(ext string) string {

	ctype := mime.TypeByExtension(ext)

	if ctype != "" {
		// Found by Go utility.
		return ctype
	}

	if ext == ".min.css" || ext == ".min.css.map" {
		return "text/css; charset=utf-8"

	} else if ext == ".js.map" || ext == ".min.js" {
		return "application/javascript"

	} else if ext == ".min.js.map" {
		// application/octet-stream works best for this, although
		// You could return application/javascript so that the content would be
		// viewable in a browser, however, while visible as text it may still
		// not work with some browsers (you may get an error in the console).
		return "application/octet-stream"
	}

	return ctype
}

// ServeStaticFile processes static files for a website. Static files
// are the ones that require no additional rending before their content
// is written to a ResponseWrite object, hence no custom error handling, if file is not found.
// The MIME is written to the Response Header according to the extension of the
// requested file. Examples are: .js, .css, .html.
func (h *HTTP) ServeStaticFile(w http.ResponseWriter, r *http.Request) {

	uriPath := r.URL.Path

	// Note about Security:
	// If you need to apply security for your static files (i.e restrict access to some .js or image files),
	// add your rules here to catch matches by path, ip addr, header, http method, etc.
	// For example, you may choose a range of IP addresses not to be able to use a
	// certian js file...You could warn the user in your API or website and then
	// make certain that your page does not leave your server.
	// The following is a crude example:
	// blockedList := []string{"###.29.29.3", "###.29.29.4", "###.29.29.5"}
	// ip := parseIPAddress(r)
	// for i := 0; i < len(blockedList); i++ {
	// 	if ip == blockedList[i] {
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	// 		return
	// 	}
	// }

	ext := getFileExtension(uriPath)

	if !strings.Contains(cachetypes, ext) {
		// This is web page like .html .pl,... that is cached by this method.
		rPath := r.URL.Path
		physPath := fmt.Sprintf("%s%s", h.RootPhysicalPath, rPath)

		http.ServeFile(w, r, physPath)
		return
	}

	h.setContentTypeAndWrite(w, r)
}

// AddSuffix adds file extension (i.e. .html) to the path if not present.
// It will check for /null in the path (maybe passed by javascript in error).
// It also adds index.html to the path, if the path is a directory.
func (h *HTTP) AddSuffix(rPath string, fileExtension string) string {

	if rPath == "/null" {
		return "/"
	}

	if strings.HasSuffix(rPath, fileExtension) {
		return rPath
	}

	physPath := fmt.Sprintf("%s%s%s", h.RootPhysicalPath, rPath, fileExtension)
	if fileOrDirectoryExists(physPath) {
		rPath = fmt.Sprintf("%s%s", rPath, fileExtension)
		return rPath
	}

	// If a directory - add index.html
	newRpath := fmt.Sprintf("%s/index%s", rPath, fileExtension)
	physPath = fmt.Sprintf("%s%s", h.RootPhysicalPath, newRpath)
	if fileOrDirectoryExists(physPath) {
		return newRpath
	}

	// TODO:
	// Add your customized handler here. For example you may want
	// to see if your path is a few directires deep... i.e.
	// http://mywebsite/mydir1/mydir2/mysubjectdir and then add the
	// /index.html to the above path.
	// You could also choose to accept URL paths with or without the
	// .html.

	// as-is
	return rPath
}
