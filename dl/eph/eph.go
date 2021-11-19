package eph

import (
	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

// implemented by individual commands
type Ephemera interface {
	Catalog(c *Catalog, d *Domain, at string) error
	Phase() Phase
}

type Phase int

//go:generate stringer -type=Phase
const (
	Domains Phase = iota
	Plurals
	Ancestry
	Fields
	Relations
	Defaults
	Nouns
	Relatives
	Patterns
	Grammar
	Tests
	References
	//
	NumPhases
)

// words in Tapestry are "normalized" for easier comparison.
// whitespace is collapsed and replaced with single underscores.
// punctuation gets removed entirely.
// letters are lowercased.
func UniformString(s string) (ret string, okay bool) {
	out := lang.LowerBreakcase(s)
	return out, len(out) > 0
}

func InvalidString(str string) error {
	return invalidStringError{str}
}

type invalidStringError struct {
	str string
}

func (x invalidStringError) Error() string {
	return errutil.Sprint("invalid string %q", x.str)
}
