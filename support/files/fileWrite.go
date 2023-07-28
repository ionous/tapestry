package files

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"syscall"
)

func WriteJson(outPath string, data interface{}, pretty bool) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = WriteJsonFile(fp, data, pretty)
	}
	return
}

func WriteJsonFile(fp *os.File, data interface{}, pretty bool) (err error) {
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
	return
}

//

const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
)

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
