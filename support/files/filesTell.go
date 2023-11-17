package files

import (
	"os"

	"github.com/ionous/tell"
)

// serialize to the passed path
func WriteTell(outPath string, data any) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = WriteTellFile(fp, data)
	}
	return
}

// serialize to the passed open file
func WriteTellFile(fp *os.File, data any) (err error) {
	enc := tell.NewEncoder(fp)
	return enc.Encode(data)
}
