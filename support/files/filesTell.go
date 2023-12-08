package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/tell"
	"github.com/ionous/tell/decode"
	"github.com/ionous/tell/encode"
	"github.com/ionous/tell/maps"
	"github.com/ionous/tell/notes"
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

// serialize to the passed open file
func WriteTellFile(fp *os.File, data any) (err error) {
	enc := tell.NewEncoder(fp)
	var m encode.MapTransform
	enc.SetMapper(m.
		// keys with -- are metadata
		KeyTransform(func(key r.Value) string {
			str := key.String()
			if strings.HasPrefix(str, "--") {
				str = str[2:]
			}
			return str
		}).
		CommentFactory(tapComments).
		Mapper())
	//
	var n encode.SequenceTransform
	enc.SetSequencer(n.
		// tap only stores comments in its commands
		CommentLocation(encode.NoComments).
		Sequencer())
	//
	return enc.Encode(data)
}

func ReadTellFile(fp *os.File, pv *map[string]any) (err error) {
	dec := decode.MakeDecoder(
		nil,
		notes.DiscardComments(),
	)
	dec.SetMapper(func(reserve bool) maps.Builder {
		tap := make(tapMap)
		x, y := dec.Position()
		tap["--pos"] = fmt.Sprintf("%d,%d", y, x)
		return tap
	})
	if raw, e := dec.Decode(bufio.NewReader(fp)); e != nil {
		err = e
	} else if raw == nil {
		*pv = nil
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
	return
}

type tapMap map[string]any

func (m tapMap) Add(key string, val any) maps.Builder {
	if len(key) > 0 && unicode.IsLower(rune(key[0])) {
		key = "--" + key
	}
	m[key] = val
	return m
}

func (m tapMap) Map() any {
	return (map[string]any)(m)
}

func tapComments(v r.Value) (ret encode.CommentIter, err error) {
	if k := v.Kind(); k != r.Interface {
		err = fmt.Errorf("expected an interface value; got %s(%s)", k, v.Type())
	} else {
		switch val := v.Elem(); {
		case val.Kind() == r.String:
			ret = newComment("# " + val.String())
		case val.Kind() == r.Slice:
			var join strings.Builder
			for i, cnt := 0, val.Len(); i < cnt; i++ {
				at := val.Index(i)
				if k := at.Kind(); k != r.Interface {
					err = errutil.Fmt("comment slice contains %s not interface", k)
				} else if at := at.Elem(); at.Kind() != r.String {
					err = errutil.New("comment slice contains an underlying %s", at.Kind())
				} else {
					if i == 0 {
						join.WriteString("# ")
					} else {
						join.WriteString("\n\t# ")
					}
				}
				str := at.String()
				join.WriteString(str)
			}
			ret = newComment(join.String())

		default:
			err = errutil.Fmt("unexpected kind (%s) of comment", val.Kind())
		}
	}
	return
}

func newComment(str string) *tapComment {
	return &tapComment{next: str}
}

type tapComment struct {
	text, next string
}

func (s *tapComment) Next() bool {
	s.text, s.next = s.next, ""
	return len(s.text) > 0
}
func (s *tapComment) GetComment() encode.Comment {
	return encode.Comment{} /// fix: s.text
}
