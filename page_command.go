package webutil

import (
	"bytes"
	"fmt"
)

// parsePageCommand processes command directives inside a page.
func (h *HTTP) parsePageCommand(b []byte, cmd string) ([]byte, error) {

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

			relPath := bytes.Split(bytes.Split(v[i], delmByte)[0], []byte(":"))[1]
			pysPath := fmt.Sprintf("%s/%s", h.RootPhysicalPath, relPath)
			bx, err := ReadFile(pysPath)

			if err != nil {
				return nil, err
			}

			phrase := fmt.Sprintf("%s:%s%s", string(leftSearchLeafe), relPath, delm)

			b = bytes.ReplaceAll(b, []byte(phrase), bx)
		}
	}

	return b, nil
}

// ProcessPageCommands replaces command directive blocks with
// their results. In the following example the content of
// the file /web/html/index.html will be placed inside the
// div tag
//   <div style="border:none">
//     {{.$LoadFile:/web/html/index.html}}
//   </div>
func (h *HTTP) ProcessPageCommands(b []byte) ([]byte, error) {

	b, err := h.parsePageCommand(b, PageCmdLoadFile)

	return b, err
}
