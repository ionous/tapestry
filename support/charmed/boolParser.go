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

func (p *BoolParser) String() string {
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
		ret = charm.RunState(r, MatchString("true", func() {
			p.result = 1
		}))
	case 'f':
		ret = charm.RunState(r, MatchString("false", func() {
			p.result = -1
		}))
	}
	return
}

// match the string and call the passed function when matched
// returns error if mismatched
func MatchString(str string, matched func()) charm.State {
	var i int // index in str
	return charm.Self("match "+str, func(self charm.State, r rune) (ret charm.State) {
		if i >= len(str) {
			return nil // really should never get here unless the string is empty
		} else if match, size := utf8.DecodeRuneInString(str[i:]); match != r {
			e := errutil.New("mismatched string")
			ret = charm.Error(e)
		} else if i += size; i < len(str) {
			ret = self // loop
		} else {
			matched()
			ret = charm.Finished("bool")
		}
		return
	})
}
