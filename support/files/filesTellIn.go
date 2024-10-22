package files

import (
	"bufio"
	"io"
	"os"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"github.com/ionous/tell/collect"
	"github.com/ionous/tell/collect/stdmap"
	"github.com/ionous/tell/collect/stdseq"
	"github.com/ionous/tell/decode"
	"github.com/ionous/tell/note"
)

type Ofs struct {
	File string
	Line int
}

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

func ReadTell(in io.Reader) (any, error) {
	return ReadTellRunes(bufio.NewReader(in), Ofs{}, true)
}

// reads until the passed reader is exhausted ( hits eof )
// returns nil error when finished.
// includeComments helps with testing
func ReadTellRunes(in io.RuneReader, ofs Ofs, includeComments bool) (ret any, err error) {
	// fix: the decoder should should take a line offset
	// rather than the mucking about done in here.
	dec := decode.Decoder{UseFloats: true} // sadly, that's all tapestry supports. darn json.
	if !includeComments {
		dec.SetMapper(stdmap.Make)
		dec.SetSequencer(stdseq.Make)
	} else {
		var docComments note.Book
		dec.SetMapper(func(reserve bool) collect.MapWriter {
			m := make(tapMap)
			x, y := dec.Position()
			m[compact.Position] = []int{x, y + ofs.Line}
			if len(ofs.File) > 0 {
				m[compact.File] = ofs.File
			}
			return m
		})
		dec.SetSequencer(func(reserve bool) collect.SequenceWriter {
			return make(tapSeq, 0)
		})
		dec.UseNotes(&docComments)
	}
	if v, e := dec.Decode(in); e == nil {
		ret = v
	} else if pos, ok := e.(decode.ErrorPos); !ok {
		err = e
	} else {
		y, x := pos.Pos()
		err = decode.ErrorAt(y+ofs.Line, x, pos.Unwrap())
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
		if !compact.IsMarkup(key) {
			if val == nil {
				// replace unary values
				key = key[:len(key)-1]
				val = true
			}
		} else {
			if end := len(key) - 1; key[end] == ':' {
				key = key[:end]
			}
		}
		m[key] = val
	} else {
		// tbd: would it make more sense to send around "Comment" structs?
		if str := val.(string); len(str) > 0 {
			m[compact.Comment] = readComment(str)
		}
	}
	return m
}

// split into separate lines and remove the leading hashes.
// the returned data is either a string, or a plain data slice of strings
// ( ie. any[]{"example"} )
func readComment(str string) (ret any) {
	var last int
	var lines []any

	if i := strings.IndexRune(str, '#'); i >= 0 {
		// ignore any leading inline markers
		str = str[i:]

		// search for newlines
		for i, ch := range str {
			if ch == '\n' {
				line := chopHash(str[last:i])
				lines = append(lines, line)
				last = i + 1 // skip the newline
			}
		}
		//
		if last == 0 {
			// only ever one line? return a string.
			ret = chopHash(str)
		} else {
			// chop any trailing comment that didnt end in a newline
			if rest := str[last:]; len(rest) > 0 {
				line := chopHash(rest)
				ret = append(lines, line)
			} else {
				ret = lines
			}
		}
	}
	return
}

func chopHash(str string) (ret string) {
	if cnt := len(str); cnt == 0 || str[0] != '#' {
		panic("chopHash expected a comment line")
	} else if cnt > 1 {
		// just the hash means an empty comment line
		// otherwise we assume a hash and a space
		ret = str[2:]
	}
	return
}
