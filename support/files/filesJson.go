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
