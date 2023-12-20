package files

import (
	"fmt"
	"os"
	"strings"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/tell"
	"github.com/ionous/tell/encode"
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
	var n encode.SequenceTransform
	enc.SetSequencer(n.
		// tap only stores comments in its commands
		CommentLocation(encode.NoComments).
		Sequencer())
	enc.SetMapper(makeMapping)
	return enc.Encode(data)
}

func makeMapping(src r.Value) (ret encode.MappingIter, err error) {
	if src.Len() == 0 {
		err = errutil.New("can't encode empty command")
	} else {
		var pairs []pair
		var header []string
		for it := src.MapRange(); it.Next(); {
			k, v := it.Key().String(), it.Value()
			if !strings.HasPrefix(k, "--") {
				if cnt := len(k); cnt > 0 && k[cnt-1] != ':' {
					k, v = k+":", r.ValueOf(nil)
				}
				pairs = append([]pair{{k, v}}, pairs...)
			} else {
				if len(k) > 2 {
					k = k[2:] + ":"
					pairs = append(pairs, pair{k, v})
				} else if c, e := encodeComments(v); e != nil {
					err = e
					break
				} else {
					header = c
				}
			}
		}
		if err == nil {
			ret = &mapIter{pairs: pairs, header: header}
		}
	}
	return
}

type pair struct {
	key string
	val r.Value
}

type mapIter struct {
	pair    pair
	comment encode.Comment
	pairs   []pair
	header  []string
}

func (m *mapIter) Next() (okay bool) {
	if okay = len(m.pairs) > 0; okay {
		m.pair, m.pairs = m.pairs[0], m.pairs[1:]
		m.comment, m.header = encode.Comment{Header: m.header}, nil
	}
	return
}

func (m *mapIter) GetKey() string             { return m.pair.key }
func (m *mapIter) GetValue() any              { return m.pair.val.Interface() }
func (m *mapIter) GetReflectedValue() r.Value { return m.pair.val }
func (m *mapIter) GetComment() encode.Comment { return m.comment }

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
						err = errutil.Fmt("comment slice contains %s not interface", k)
						break
					} else if at := at.Elem(); at.Kind() != r.String {
						err = errutil.New("comment slice contains an underlying %s", at.Kind())
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
			err = errutil.Fmt("unexpected kind (%s) of comment", val.Kind())
		}
	}
	return
}
