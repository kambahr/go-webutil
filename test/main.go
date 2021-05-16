package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/kambahr/go-webutil"
)

var mWebutil *webutil.HTTP
var mInstallPath string

// This is a simple website to show the usage of the webuti package.

func main() {

	var portNo int = 8005

	// This is the default root directory, where your main is
	// executed from; but it could be any path on your system -- that
	// you've located your web-files.
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	mInstallPath = dir

	// This covers all supporting files like js, css, img,...
	mWebutil = webutil.NewHTTP(mInstallPath, 5*time.Minute)
	http.HandleFunc("/assets/", mWebutil.ServeStaticFile)

	http.HandleFunc("/", handleMyPage)

	svr := http.Server{
		Addr:           fmt.Sprintf(":%d", portNo),
		MaxHeaderBytes: 20480,
	}

	fmt.Printf("Listening to port %d\n", portNo)

	log.Fatal(svr.ListenAndServe())
}

func handleMyPage(w http.ResponseWriter, r *http.Request) {

	var b []byte

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Relative path to store on the cache list.
	const uri string = "/mypage.html"
	rPath := strings.ToLower(r.URL.Path)

	pageValid := (strings.HasSuffix(rPath, "/") || strings.HasSuffix(rPath, uri))

	if !pageValid {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h3>Error 404 - Not Found</h3>")
		return
	}

	// Read the file bytes.
	tInshtmlPath := fmt.Sprintf("%s%s", mInstallPath, uri)

	fExistOnDisk := fileExists(tInshtmlPath)

	if !fExistOnDisk {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h3>Error 404 - Not Found</h3>")
		return
	}

	b, _ = ioutil.ReadFile(tInshtmlPath)

	// Write the response.
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
