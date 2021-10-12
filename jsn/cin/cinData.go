package cin

import (
	"encoding/json"
	"unicode"
)

type cinMap map[string]json.RawMessage

const eof = rune(-1)

type sigReader struct {
	cmd    string
	params []string
	buf    []rune
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
