package files

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"unicode"

	r "reflect"

	"github.com/ionous/tell/collect"
	"github.com/ionous/tell/decode"
	"github.com/ionous/tell/note"
)

// deserialize from the passed path
func ReadTell(inPath string, pv *map[string]any) (err error) {
	if fp, e := os.Open(inPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = ReadTellFile(fp, pv)
	}
	return
}

func ReadTellFile(fp fs.File, pv *map[string]any) (err error) {
	var docComments note.Book
	dec := decode.Decoder{UseFloats: true} // sadly, that's all tapestry supports. darn json.
	dec.SetMapper(func(reserve bool) collect.MapWriter {
		return make(tapMap)
	})
	dec.SetSequencer(func(reserve bool) collect.SequenceWriter {
		return make(tapSeq, 0, 0)
	})
	dec.UseNotes(&docComments)
	if raw, e := dec.Decode(bufio.NewReader(fp)); e != nil {
		err = e
	} else {
		if raw == nil {
			*pv = make(map[string]any)
		} else {
			out, res := r.ValueOf(pv).Elem(), r.ValueOf(raw)
			if rt, ot := res.Type(), out.Type(); rt.AssignableTo(ot) {
				out.Set(res)
			} else if res.CanConvert(ot) {
				out.Set(res.Convert(ot))
			} else {
				err = fmt.Errorf("result of %q cant be written to a pointer of %q", rt, ot)
			}
		}
	}
	return
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
			lines := cleanComment(str)
			m["--"] = packComment(lines)
		}
	}
	return m
}

func packComment(lines []string) (ret any) {
	if len(lines) == 1 {
		ret = lines[0]
	} else {
		ret = lines
	}
	return
}

// split into separate lines, chop the leading hashes, and avoid leading newlines
// ( that last gets used to indicate a not "inline" comment )
func cleanComment(str string) []string {
	lines := strings.Split(str, "\n")
	if first := lines[0]; len(first) == 0 {
		lines = lines[1:]
	}
	for i, el := range lines {
		lines[i] = el[2:]
	}
	return lines
}