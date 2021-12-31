# Web Utility for Golang websites.

## Webutil is a simple, lightweight utility for Golang websites.

It loads, selects the right mime, and caches static (js, css, image,...), and html files.
It also has a place-holder to apply security in order to withhold a file (i.e. a javascript file).

## Features

### MIME Types
MIME types are configured and the corrected response header is written by default.

### HTTP Response Wrapper

#### WriteResponse
Compression headers are added with each response; just pass the data that is to be displayed on a page
to the WriteResponse function:

```go
WriteResponse(data []byte, w http.ResponseWriter, r *http.Request)
```

#### HTTPExec
HTTPExec is an http wrapper. You can set timeout, pass headers, and have an option of receiving
a log of the call. It returns the entire response for further parsing. 
```go
HTTPExec(method HTTPMethod, urlx string, hd http.Header, data []byte, tMillisec uint, logRequest bool) (HTTPResult, error)
```

#### Serve Static Files
ServeStaticFile can be setup for processing supporting files like css, js, and image files. 
```go
func (h *HTTP) ServeStaticFile(w http.ResponseWriter, r *http.Request)
````
All needed is assigning the path to the http handler as the following example:
````go
mWebutil := webutil.NewHTTP(mInstallPath, 5*time.Minute)
http.HandleFunc("/assets/", mWebutil.ServeStaticFile)
````
#### Process Page Directives
Insert page directives inside html pages. 
```go
func (h *HTTP) ProcessPageCommands(b []byte) ([]byte, error)
````
Use the **LoadFile** directive to insert content inside a block; easily reuse pieces of code inside pages without writing separate code.
Usage Example:
```
<div style="border:none">
  {{.$LoadFile:/web/html/my-cool-grid.html}}
</div>
```

### Comments
You can leave comments in any file (.html, .js, .go,..), knowing that they will not reach the client.
````
{{.COMMENTS <your comments go here> }}
````

#### Run the test app

- Start a shell window in sample directory.
- go build -o webutilDemo && ./webutilDemo
- Navigate to http://localhost:8005/mypage.html
