package charmed

import (
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// parses `true` or `false`
type BoolParser struct {
	result int
}

func (p *BoolParser) StateName() string {
	return "bools"
}

func (p *BoolParser) GetBool() (ret bool, err error) {
	switch p.result {
	case 1:
		ret = true
	case -1:
		ret = false
	default:
		err = errutil.New("unmatched bool")
	}
	return
}

func (p *BoolParser) NewRune(r rune) (ret charm.State) {
	switch r {
	case 't':
		ret = charm.RunStep(r, MatchString("true"), charm.OnExit("matched", func() {
			p.result = 1
		}))
	case 'f':
		ret = charm.RunStep(r, MatchString("false"), charm.OnExit("matched", func() {
			p.result = -1
		}))
	}
	return
}

func MatchString(str string) charm.State {
	var i int // index in str
	return charm.Self("match "+str, func(self charm.State, r rune) (ret charm.State) {
		if i < len(str) { // returns nil once we've matched the whole string
			if match, size := utf8.DecodeRuneInString(str[i:]); match == r {
				ret = self
				i += size
			} else {
				e := errutil.New("mismatched string")
				ret = charm.Error(e)
			}
		}
		return
	})
}
