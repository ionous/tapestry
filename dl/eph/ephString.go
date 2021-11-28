package eph

import (
	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

// words in Tapestry are "normalized" for easier comparison.
// whitespace is collapsed and replaced with single underscores.
// punctuation gets removed entirely.
// letters are lowercased.
func UniformString(s string) (ret string, okay bool) {
	out := lang.Underscore(s)
	return out, len(out) > 0
}

func UniformStrings(strs []string) (ret []string, err error) {
	for _, src := range strs {
		if s, ok := UniformString(src); !ok {
			err = errutil.Append(err, InvalidString(src))
		} else {
			ret = append(ret, s)
		}
	}
	return
}

func InvalidString(str string) error {
	return invalidStringError{str}
}

type invalidStringError struct {
	str string
}

func (x invalidStringError) Error() string {
	return errutil.Sprintf("invalid string %q", x.str)
}
