package cmdcompact

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

type format int

const (
	jsonFormat format = iota
	tellFormat
)

func (f format) write(out string, data any, pretty bool) (err error) {
	switch f {
	case jsonFormat:
		err = files.WriteJson(out, data, pretty)
	case tellFormat:
		err = files.WriteTell(out, data)
	default:
		err = errutil.New("unknown format")
	}
	return
}

func (f format) read(in string) (ret map[string]any, err error) {
	switch f {
	case jsonFormat:
		if b, e := files.ReadFile(in); e != nil {
			err = e
		} else {
			err = json.Unmarshal(b, &ret)
		}
	case tellFormat:
		err = files.ReadTell(in, &ret)
	default:
		err = errutil.New("unknown format")
	}
	return
}
