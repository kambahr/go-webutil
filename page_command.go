package webutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// parseJsonIntoSingleLine trimes json string into a single line.
func (h *HTTP) parseJsonIntoSingleLine(bx []byte) []byte {

	m := make([]map[string]interface{}, 0)
	json.Unmarshal(bx, &m)

	b, _ := json.Marshal(m)
	//b = bytes.ReplaceAll(b, []byte(`\"`), []byte(`"`))
	xs := string(b)

	return b

	// var emptyByte []byte
	// bx = bytes.TrimSpace(bytes.ReplaceAll(bx, []byte("\n"), emptyByte))
	// bx = bytes.TrimSpace(bytes.ReplaceAll(bx, []byte("\t"), emptyByte))
	// bx = bytes.TrimSpace(bytes.ReplaceAll(bx, []byte("\r"), emptyByte))
	// bx = bytes.ReplaceAll(bx, []byte(`\"`), []byte(`"`))

	// tx := [][2]string{
	// 	{",", `"`},
	// 	{"[", "{"},
	// 	{"{", `"`},
	// }

	// for i := 0; i < len(tx); i++ {
	// 	r := regexp.MustCompile(fmt.Sprintf(`\%s(.*?)\%s`, tx[i][0], tx[i][1]))
	// 	m := r.FindAllString(string(bx), -1)
	// 	sx := fmt.Sprintf("%s%s", tx[i][0], tx[i][1])
	// 	for i := 0; i < len(m); i++ {
	// 		mTrimed := strings.ReplaceAll(string(bx), m[i], sx)
	// 		bx = []byte(mTrimed)
	// 	}
	// }

	// return bx
}

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
			bx = h.parseJsonIntoSingleLine(bx)

			phrase := fmt.Sprintf("%s:%s%s", string(leftSearchLeafe), relPath, delm)

			b = bytes.ReplaceAll(b, []byte(phrase), bx)

		}
	} else if cmd == "$RenderJSONToHTMLControl" {

		v := strings.Split(string(b), string(leftSearchLeafe))

		for i := 0; i < len(v); i++ {
			v[i] = strings.TrimSpace(strings.ReplaceAll(v[i], "\n", ""))
			v[i] = strings.TrimSpace(strings.ReplaceAll(v[i], "\r", ""))
			v[i] = strings.TrimSpace(strings.ReplaceAll(v[i], "\t", ""))

			if !strings.Contains(v[i], delm) {
				continue
			}
			if !strings.HasPrefix(v[i], `:[{"`) {
				continue
			}
			jsn := v[i]
			jsn = jsn[1:]                   // Minus :
			pos := strings.Index(jsn, "}}") // find the end
			jsn = jsn[:pos]
			if len(jsn) < 5 {
				continue
			}

			bJsn := []byte(jsn)
			bJsn = h.parseJsonIntoSingleLine(bJsn)
			m := make([]map[string]interface{}, 0)
			json.Unmarshal(bJsn, &m)

			// Expect bootstrap to be installed

			html := ""
			hcheck := `
<div class="form-check mb-2 mt-2">
   <input class="form-check-input" type="checkbox" value="" id="#CHECK_ID#" #CHECKED#>
	<label class="form-check-label" for="#CHECK_ID#">#LABEL_VALUE#</label>
</div>`
			h := `
<div class="input-group mb-3">
  <span class="input-group-text" id="#LABEL_ID#">#LABEL#</span>
  <input type="text" class="form-control" value="#VALUE#" id="#INPT_ID#" placeholder="#LABEL#" aria-label="#LABEL#" aria-describedby="#LABEL_ID#">
</div>`
			for _, v := range m {
				for kk, vv := range v {
					key := kk
					if strings.Contains(key, "-") {
						// remove spaces
						key = strings.Title(strings.ReplaceAll(key, " ", ""))
					} else {
						// caps first letter
						key = strings.Title(key)

						// and then remove space
						key = strings.Title(strings.ReplaceAll(key, " ", ""))
					}
					if reflect.TypeOf(vv).Kind() == reflect.String {
						lbl := strings.Title(strings.ReplaceAll(key, "-", " "))

						cx := cases.Title(language.English)
						lbl = cx.String(lbl)

						lblID := strings.ToLower(fmt.Sprintf("label-%s", key))
						inptID := strings.ToLower(fmt.Sprintf("input-%s", key))
						value := fmt.Sprintf("%v", vv)
						hx := strings.ReplaceAll(h, "#LABEL#", lbl)
						hx = strings.ReplaceAll(hx, "#VALUE#", value)
						hx = strings.ReplaceAll(hx, "#LABEL_ID#", lblID)
						hx = strings.ReplaceAll(hx, "#INPT_ID#", inptID)
						html = fmt.Sprintf("%s%s", html, hx)

					} else if reflect.TypeOf(vv).Kind() == reflect.Slice {
						ifarry := vv.([]interface{})

						isString := reflect.TypeOf(ifarry[0]).Kind() == reflect.String

						kk = strings.Title(strings.ReplaceAll(kk, "-", " "))

						cx := cases.Title(language.English)
						kk = cx.String(kk)

						html = fmt.Sprintf(`%s<dl><dt class="mb-1">%s</dt>`, html, kk)
						if isString {
							lst := `<ul class="list-group" style="border:none">`
							for j := 0; j < len(ifarry); j++ {
								lst = fmt.Sprintf(`%s<dd><input type="text" class="form-control" value="%s"></dd>`, lst, ifarry[j])
							}
							html = fmt.Sprintf(`%s%s`, html, lst)
						}
						html = fmt.Sprintf(`%s</dl>`, html)
					} else if reflect.TypeOf(vv).Kind() == reflect.Bool {
						chkVal := strings.ToLower(fmt.Sprintf("%v", vv))

						cx := cases.Title(language.English)
						chkLbl := cx.String(kk)

						chk := strings.ReplaceAll(hcheck, "#LABEL_VALUE#", chkLbl)
						chkID := fmt.Sprintf("chk-%s", strings.ToLower(strings.Title(strings.ReplaceAll(key, "-", " "))))
						chk = strings.ReplaceAll(chk, "#CHECK_ID#", chkID)

						if chkVal == "true" {
							chk = strings.ReplaceAll(chk, "#LABEL_VALUE#", "checked")
						} else {
							chk = strings.ReplaceAll(chk, "#LABEL_VALUE#", "")
						}
						html = fmt.Sprintf(`%s%s`, html, chk)
					}
				}
			}

			toReplace := fmt.Sprintf("%s:%s}}", string(leftSearchLeafe), jsn)

			b = bytes.ReplaceAll(b, []byte(toReplace), []byte(html))
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

	var err []error
	var cmdArry = []string{
		PageCmdLoadFile,
		"$RenderJSONToHTMLControl",
	}

	for i := 0; i < len(cmdArry); i++ {
		b, err = h.parsePageCommand(b, cmdArry[i])
	}

	return b, err
}
