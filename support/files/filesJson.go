package files

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

// serialize to the passed path
func WriteJson(outPath string, data any, pretty bool) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = WriteJsonFile(fp, data, pretty)
	}
	return
}

// serialize to the passed open file
func WriteJsonFile(fp *os.File, data any, pretty bool) (err error) {
	js := json.NewEncoder(fp)
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
