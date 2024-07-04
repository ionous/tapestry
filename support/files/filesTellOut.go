package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"github.com/ionous/tell/encode"
)

// serialize to the passed path
func SaveTell(outPath string, data any) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		err = WriteTell(fp, data)
	}
	return
}

// serialize to the passed open file
func WriteTell(w io.Writer, data any) (err error) {
	enc := TellEncoder(w)
	return enc.Encode(data)
}

func makeMapping(src r.Value) (ret encode.Iterator, err error) {
	if src.Len() == 0 {
		err = errors.New("can't encode empty command")
	} else {
		var pairs []pair
		for it := src.MapRange(); it.Next(); {
			k, v := it.Key().String(), it.Value()
			if !compact.IsMarkup(k) {
				// if there's a key that doesnt end with a add colon:
				// add one. these are unary commands.
				if cnt := len(k); cnt > 0 && k[cnt-1] != ':' {
					k, v = k+":", r.ValueOf(nil)
				}
				pairs = append(pairs, pair{k, v})
			} else {
				// prepend metadata.
				p := pair{key: k, val: v}
				pairs = append([]pair{p}, pairs...)
			}
		}
		ret = &mapIter{pairs: pairs}
	}
	return
}

type pair struct {
	key string
	val r.Value
}

type mapIter struct {
	pair  pair
	pairs []pair
}

func (m *mapIter) Next() (okay bool) {
	if okay = len(m.pairs) > 0; okay {
		m.pair, m.pairs = m.pairs[0], m.pairs[1:]
	}
	return
}

func (m *mapIter) GetKey() string             { return m.pair.key }
func (m *mapIter) GetValue() any              { return m.pair.val.Interface() }
func (m *mapIter) GetReflectedValue() r.Value { return m.pair.val }

type headerIt struct {
	header  []string
	comment encode.Comment
}

func (m *headerIt) Next() (okay bool) {
	if m.header != nil {
		m.comment, m.header = encode.Comment{Header: m.header}, nil
		okay = true
	}
	return
}
func (m *headerIt) GetComment() encode.Comment {
	return m.comment
}

func makeComments(v r.Value) (ret encode.Comments, err error) {
	if str, e := encodeComments(v); e != nil {
		err = e
	} else if len(str) > 0 {
		ret = &headerIt{header: str}
	}
	return
}

// commands have (at most) one header paragraph
func encodeComments(v r.Value) (ret []string, err error) {
	if k := v.Kind(); k != r.Interface {
		err = fmt.Errorf("expected an interface value; got %s(%s)", k, v.Type())
	} else {
		switch val := v.Elem(); {
		case val.Kind() == r.String:
			ret = []string{"# " + val.String()}

		case val.Kind() == r.Slice:
			if cnt := val.Len(); cnt > 0 {
				header := make([]string, cnt)
				for i := 0; i < cnt; i++ {
					at := val.Index(i)
					if k := at.Kind(); k != r.Interface {
						err = fmt.Errorf("comment slice contains %s not interface", k)
						break
					} else if at := at.Elem(); at.Kind() != r.String {
						err = fmt.Errorf("comment slice contains an underlying %s", at.Kind())
						break
					} else { // fix: maybe it'd make more sense for tell to do this?
						header[i] = "# " + at.String()
					}
				}
				if err == nil {
					ret = header
				}
			}
		default:
			err = fmt.Errorf("unexpected kind %q of comment", val.Kind())
		}
	}
	return
}
