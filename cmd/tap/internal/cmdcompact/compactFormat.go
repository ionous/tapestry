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
		tempCommentHack(data)
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
