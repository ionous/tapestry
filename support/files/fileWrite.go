package files

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

func WriteJson(outPath string, data interface{}, pretty bool) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		if str, ok := data.(string); ok {
			_, err = fp.Write(prettify(str, pretty))
		} else {
			js := json.NewEncoder(fp)
			js.SetEscapeHTML(false)
			if pretty {
				js.SetIndent("", "  ")
			}
			err = js.Encode(data)
		}
	}
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
