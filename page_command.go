package webutil

import (
	"bytes"
	"fmt"
)

// parsePageCommand processes command directives inside a page.
func (h *HTTP) parsePageCommand(b []byte, cmd string) ([]byte, []error) {

	var errArry []error

	delm := "}}"
	delmByte := []byte(delm)
	leftSearchLeafe := []byte(fmt.Sprintf("{{.%s", cmd))

	v := bytes.Split(b, leftSearchLeafe)

	if cmd == PageCmdLoadFile {

		// Example:
		//   <div style="border:none">
		//     {{.$LoadFile:/web/html/index.html}}
		//   </div>

		// Loop thru the page and process all occurance of {{.$LoadFile
		for i := 0; i < len(v); i++ {
			if !bytes.Contains(v[i], delmByte) {
				continue
			}
			v := bytes.Split(v[i], delmByte)
			if len(v) < 2 {
				continue
			}
			v = bytes.Split(v[0], []byte(":"))
			if len(v) < 2 {
				continue
			}
			relPath := v[1]
			pysPath := fmt.Sprintf("%s/%s", h.RootPhysicalPath, relPath)
			bx, err := ReadFile(pysPath)

			if err != nil {
				errArry = append(errArry, err)
				continue
			}

			phrase := fmt.Sprintf("%s:%s%s", string(leftSearchLeafe), relPath, delm)

			b = bytes.ReplaceAll(b, []byte(phrase), bx)
		}
	}

	return b, errArry
}

// ProcessPageCommands replaces command directive blocks with
// their results. In the following example the content of
// the file /web/html/index.html will be placed inside the
// div tag
//   <div style="border:none">
//     {{.$LoadFile:/web/html/index.html}}
//   </div>
func (h *HTTP) ProcessPageCommands(b []byte) ([]byte, []error) {

	b, err := h.parsePageCommand(b, PageCmdLoadFile)

	return b, err
}
