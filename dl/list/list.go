package list

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

// can v be inserted into els?
func IsInsertable(v, els rt.Value) (okay bool) {
	return isInsertable(els, v.Affinity(), v.Type())
}

// can v be appended to els?
// this is similar to IsInsertable, except that v can itself be a list.
func IsAppendable(v, els rt.Value) (okay bool) {
	inAff := v.Affinity()
	if unlist := affine.Element(inAff); len(unlist) > 0 {
		inAff = unlist
	}
	return isInsertable(els, inAff, v.Type())
}

func isInsertable(els rt.Value, haveAff affine.Affinity, haveType string) (okay bool) {
	okay = true // provisionally
	listAff := els.Affinity()
	if needAff := affine.Element(listAff); len(needAff) == 0 {
		okay = false // els was not actually a list
	} else if haveAff != needAff {
		okay = false // the element affinities dont match
	} else if haveAff == affine.Record && haveType != els.Type() {
		okay = false // the record types dont match
	}
	return
}

type insertError struct {
	v, els rt.Value
}

func (e insertError) Error() string {
	return fmt.Sprintf("%s of %q isn't insertable into %s of %q",
		e.v.Affinity(), e.v.Type(),
		e.els.Affinity(), e.els.Type())
}
