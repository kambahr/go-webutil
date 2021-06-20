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

	httpWrapperDemo()

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
	fmt.Println("open http://localhost:8005 in a browser to view the static page.")

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

// httpWrapperDemo shows usage of the webutil.HTTPExec() in web.go.
func httpWrapperDemo() {
	urlx := "https://go-webcache.githubsamples.com"
	fmt.Println(strings.Repeat("*", 80))
	fmt.Println("Demo: using the webutil.HTTPExec() wrapper")
	fmt.Println(strings.Repeat("-", 43))
	fmt.Println("calling", urlx, "...")
	var myHeader http.Header
	myHeader = make(http.Header, 1)
	myHeader.Set("User-Agent", "Go client;usage demo for webutil.HTTPExec() - https://github.com/kambahr/go-webutil; ")
	myHeader.Set("X-MyHeader", "my header value")
	myData := []byte(`{"mode": "test"}`)
	timeoutMillisecond := uint(600)
	st := time.Now()
	data, resp, errx := webutil.HTTPExec(webutil.GET, urlx, myHeader, myData, timeoutMillisecond)
	elsapsed := time.Since(st)
	if errx != nil {
		log.Fatal(errx)
	}
	fmt.Println("  ", resp.Status)
	rs := string(data)
	rs = strings.ReplaceAll(rs, "\n", "")
	rs = strings.ReplaceAll(rs, "  ", "")
	fmt.Println("  ", rs[:70], "...")
	fmt.Println("elsaped time:", elsapsed)
	fmt.Println(strings.Repeat("*", 80))
	fmt.Print()
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
