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
		CommentFactory(encodeComments).
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

// given the value of the blank key of the map
// should return something to walk every key in the map
// noting that .if mappings only have one term ( the command op )
func encodeComments(v r.Value) (ret encode.CommentIter, err error) {
	if k := v.Kind(); k != r.Interface {
		err = fmt.Errorf("expected an interface value; got %s(%s)", k, v.Type())
	} else {
		switch val := v.Elem(); {
		case val.Kind() == r.String:
			if str := val.String(); len(str) > 0 {
				header := []string{"# " + str}
				ret = encode.Comments([]encode.Comment{{
					Header: header,
				}})
			}
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
					ret = encode.Comments([]encode.Comment{{
						Header: header,
					}})
				}
			}
		default:
			err = errutil.Fmt("unexpected kind (%s) of comment", val.Kind())
		}
	}
	return
}
