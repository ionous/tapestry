package files

import (
	"bufio"
	"io"
	"os"
	"unicode"

	"github.com/ionous/tell/collect"
	"github.com/ionous/tell/decode"
	"github.com/ionous/tell/note"
)

// deserialize from the passed path
func LoadTell(inPath string) (ret any, err error) {
	if fp, e := os.Open(inPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		ret, err = ReadTell(fp)
	}
	return
}

func ReadTell(in io.Reader) (ret any, err error) {
	var docComments note.Book
	dec := decode.Decoder{UseFloats: true} // sadly, that's all tapestry supports. darn json.
	dec.SetMapper(func(reserve bool) collect.MapWriter {
		m := make(tapMap)
		x, y := dec.Position()
		m["--pos"] = []int{x, y}
		return m
	})
	dec.SetSequencer(func(reserve bool) collect.SequenceWriter {
		return make(tapSeq, 0, 0)
	})
	dec.UseNotes(&docComments)
	return dec.Decode(bufio.NewReader(in))
}

// tapestry sequences never have comments; so throw out the zeroth element
type tapSeq []any

// returns []any
func (m tapSeq) GetSequence() any {
	return ([]any)(m)
}
func (m tapSeq) IndexValue(idx int, val any) (ret collect.SequenceWriter) {
	if idx > 0 {
		ret = append(m, val)
	} else {
		ret = m
	}
	return
}

// output for decoder
type tapMap map[string]any

func (m tapMap) GetMap() any {
	return (map[string]any)(m)
}

func (m tapMap) MapValue(key string, val any) collect.MapWriter {
	if len(key) != 0 {
		// lowercase keys are tapestry metadata
		if !unicode.IsLower(rune(key[0])) {
			if val == nil {
				// replace unary values
				key = key[:len(key)-1]
				val = true
			}
		} else {
			key = "--" + key
			if end := len(key) - 1; key[end] == ':' {
				key = key[:end]
			}
		}
		m[key] = val
	} else {
		// tbd: would it make more sense to send around "Comment" structs?
		if str := val.(string); len(str) > 0 {
			m["--"] = readComment(str)
		}
	}
	return m
}

// split into separate lines and chop the leading hashes.
// the returned data is either a string, or a plain data slice of strings
// ( ie. any[]{"example"} )
func readComment(str string) (ret any) {
	var last int
	var out []any
	for i, ch := range str {
		if ch == '\n' {
			add := str[last+2 : i]
			out = append(out, add)
			last = i + 1 // skip the newline
		}
	}
	if last == 0 {
		ret = str[last+2:]
	} else {
		add := str[last+2:]
		out = append(out, add)
		ret = out
	}
	return
}
