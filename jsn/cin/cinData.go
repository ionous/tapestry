package cin

import (
	"encoding/json"
	"unicode"

	"github.com/ionous/errutil"
)

// ex. {"Story:":  [...]}
// except note that literals are stored as their value,
// functions with one parameter dont use the array, and
// functions without parameters are stored as simple strings.
type cinMap map[string]json.RawMessage

const eof = rune(-1)

func makeSig(s string) (ret sigReader) {
	ret.readSig(s)
	return
}

type sigReader struct {
	cmd    string
	params []string
	buf    []rune
}

type cinFlow struct {
	params    []string
	args      []json.RawMessage
	bestIndex int
}

func newFlowData(sig sigReader, arg json.RawMessage) (ret *cinFlow, err error) {
	pn := len(sig.params)
	// we allow ( require really ) single arguments to be stored directly
	// rather than embedded in an array
	// to make it optional, we'd really need a parallel parser to attempt to interpret the argument bytes in multiple ways.
	var args []json.RawMessage
	if pn == 1 {
		args = []json.RawMessage{arg}
	} else if pn > 1 {
		err = json.Unmarshal(arg, &args)
	}
	if err == nil {
		if an := len(args); pn != an {
			err = errutil.New("mismatched params and args", pn, an)
		} else {
			ret = &cinFlow{params: sig.params, args: args}
		}
	}
	return
}

func (f *cinFlow) findArg(name string) (ret json.RawMessage) {
	if i := f.bestIndex; i < len(f.params) && f.params[i] == name {
		ret, f.bestIndex = f.args[i], i+1
	} else {
		for i, n := range f.params {
			if n == name {
				ret, f.bestIndex = f.args[i], i+1
				break // next time we're most likely on the next arg.
			}
		}
	}
	return
}

func (s *sigReader) readSig(str string) {
	for _, r := range str {
		s.readRune(r)
	}
	s.readRune(eof)
}

func (s *sigReader) readRune(r rune) {
	if len(s.cmd) == 0 {
		s.readCmd(r)
	} else {
		s.readArg(r)
	}
}

func (s *sigReader) readCmd(r rune) {
	switch {
	case r == ':':
		s.params = append(s.params, "") // blank, unlabeled
		fallthrough
	case r == ' ' || r == eof:
		s.cmd = string(s.buf) // we have a full command name now.
		s.buf = nil
	case unicode.IsUpper(r):
		s.addSep(r)
	default:
		s.buf = append(s.buf, r)
	}
}

func (s *sigReader) readArg(r rune) {
	switch {
	case r == eof:
		if len(s.buf) == 0 {
			break // in case there was no final arg.
		}
		fallthrough
	case r == ':':
		s.params = append(s.params, string(s.buf))
		s.buf = nil
	case unicode.IsUpper(r):
		s.addSep(r)
	default:
		s.buf = append(s.buf, r)
	}
}

// camelCase to break_case helper.
func (s *sigReader) addSep(r rune) {
	if l := unicode.ToLower(r); len(s.buf) > 0 {
		s.buf = append(s.buf, '_', l)
	} else {
		s.buf = append(s.buf, l)
	}
}
