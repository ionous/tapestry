package files

import (
	"fmt"
	"io"
	"os"
)

// write a .if or .tell tapestry file
// determines the type from the passed path.
func FormattedLoad(inPath string, pv *map[string]any) (err error) {
	if fp, e := os.Open(inPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = FormattedRead(fp, Ext(inPath), pv)
	}
	return
}

// write a .if or .tell tapestry file
// determines the type from the passed path.
func FormattedSave(outPath string, data any, pretty bool) (err error) {
	switch ext := Ext(outPath); {
	case ext.Json():
		err = SaveJson(outPath, data, pretty)
	case ext.Tell():
		err = SaveTell(outPath, data)
	default:
		err = fmt.Errorf("unknown format %q", ext)
	}
	return
}

// read a .if or .tell tapestry file
func FormattedRead(in io.Reader, ext Extension, pv *map[string]any) (err error) {
	switch {
	case ext.Json():
		err = ReadJson(in, pv)
	case ext.Tell():
		if d, e := ReadTell(in); e != nil {
			err = e
		} else if m, ok := d.(map[string]any); !ok {
			err = fmt.Errorf("expected a tell mapping")
		} else {
			*pv = m
		}
	default:
		err = fmt.Errorf("unknown format %q", ext)
	}
	return
}

// write a .if or .tell story file
func FormattedWrite(w io.Writer, data any, ext Extension, pretty bool) (err error) {
	switch {
	case ext.Json():
		err = JsonEncoder(w, MakeJsonFlags(pretty, false)).Encode(data)
	case ext.Tell():
		err = WriteTell(w, data)
	default:
		err = fmt.Errorf("unknown format %q", ext)
	}
	return
}
