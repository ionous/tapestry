package files

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/ionous/errutil"
)

// write a .if or .tell story file
func FormattedSave(outPath string, data any, pretty bool) (err error) {
	switch ext := Ext(outPath); {
	case ext.Json():
		err = writeJson(outPath, data, pretty)
	case ext.Tell():
		tempCommentHack(data)
		err = SaveTell(outPath, data)
	default:
		err = errutil.New("unknown format")
	}
	return
}

// write a .if or .tell story file
func FormattedWrite(w io.Writer, data any, ext Extension, pretty bool) (err error) {
	switch {
	case ext.Json():
		err = writeJsonFile(w, data, pretty)
	case ext.Tell():
		err = WriteTell(w, data)
	default:
		err = errutil.New("unknown format")
	}
	return
}

// serialize to the passed path
func writeJson(outPath string, data any, pretty bool) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = writeJsonFile(fp, data, pretty)
	}
	return
}

// serialize to the passed open file
func writeJsonFile(w io.Writer, data any, pretty bool) (err error) {
	js := json.NewEncoder(w)
	js.SetEscapeHTML(false)
	if pretty {
		js.SetIndent("", "  ")
	}
	err = js.Encode(data)
	return
}

func prettify(str string, pretty bool) (ret []byte) {
	ret = []byte(str)
	if pretty {
		var indent bytes.Buffer
		if e := json.Indent(&indent, ret, "", "  "); e != nil {
			log.Println(e)
		} else {
			ret = indent.Bytes()
		}
	}
	return
}

// change stand alone comments "--"; embedding them into the next element of an array
func tempCommentHack(data any) {
	m := data.(map[string]any)
	for k, v := range m {
		switch v := v.(type) {
		case []any:
			m[k] = tempHackSlice(v)
		case map[string]any:
			tempCommentHack(v)
		}
	}
}

func tempHackSlice(data []any) (ret []any) {
	var comment any
	for _, el := range data {
		switch m := el.(type) {
		default:
			ret = append(ret, el)
		case []any:
			ret = append(ret, tempHackSlice(m))
		case map[string]any:
			tempCommentHack(m)

			// doesnt have a comment entry?
			if c, ok := m["--"]; !ok {
				// add the previous comment if it existed
				if comment != nil {
					m["--"] = comment
					comment = nil
				}
				// keep the element
				ret = append(ret, el)
			} else {
				// comment only? store.
				if len(m) == 1 {
					if comment != nil {
						panic("yyy")
					}
					comment = c
				} else {
					if comment != nil && comment.(string) != "" {
						panic("zzz")
					}
					ret = append(ret, el)
					comment = nil
				}
			}
		}
	}
	return
}
