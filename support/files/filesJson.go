package files

import (
	"encoding/json"
	"io"
	"os"
)

// serialize to the passed path
func SaveJson(outPath string, data any, pretty bool) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		err = JsonEncoder(fp, MakeJsonFlags(pretty, false)).Encode(data)
		fp.Close()
	}
	return
}

// deserialize from the passed path
func LoadJson(inPath string, out any) (err error) {
	if fp, e := os.Open(inPath); e != nil {
		err = e
	} else {
		err = ReadJson(fp, out)
		fp.Close()
	}
	return
}

func ReadJson(in io.Reader, out any) (err error) {
	d := json.NewDecoder(in)
	return d.Decode(out)
}
