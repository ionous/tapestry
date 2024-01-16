package files

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ionous/errutil"
)

// write a .if or .tell story file
// determines the type from the passed path.
func FormattedSave(outPath string, data any, pretty bool) (err error) {
	switch ext := Ext(outPath); {
	case ext.Json():
		err = SaveJson(outPath, data, pretty)
	case ext.Tell():
		err = SaveTell(outPath, data)
	default:
		err = errutil.New("unknown supported format", ext)
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
		err = errutil.New("unknown supported format", ext)
	}
	return
}

// serialize to the passed path
func SaveJson(outPath string, data any, pretty bool) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = JsonEncoder(fp, MakeJsonFlags(pretty, false)).Encode(data)
	}
	return
}

// panics if the passed data isnt json friendly.
func Stringify(data any) (ret string) {
	if a, e := json.MarshalIndent(data, "", " "); e != nil {
		panic(e)
	} else {
		ret = string(a)
	}
	return
}

// remove go's trailing newline.
// https://github.com/golang/go/issues/37083
// or, they could have just provided an Set<Option> like SetEscapeHTML...
type noNewLine struct {
	out     io.Writer
	pending bool
}

func (n *noNewLine) Write(p []byte) (ret int, err error) {
	if cnt := len(p); cnt > 0 {
		if n.pending {
			n.out.Write([]byte{newline})
			n.pending = false
		}
		pending := p[cnt-1] == newline
		if pending {
			n.pending = pending
			cnt--
		}
		if cnt == 0 {
			ret = 1
		} else {
			c, e := n.out.Write(p[:cnt])
			if pending && c > 0 {
				c++
			}
			ret, err = c, e
		}
	}
	return
}

const newline byte = '\n'
