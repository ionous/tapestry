package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	g "github.com/ionous/iffy/rt/generic"
)

var Slots = []composer.Slot{{
	Name: "list_getter",
	Type: (*ListGetter)(nil),
	Desc: "List getter: Helper for accessing lists.",
}, {
	Name: "list_sorter",
	Type: (*Sorter)(nil),
	Desc: "List sorter: Helper for sorting lists.",
}}

var Slats = []composer.Composer{
	(*At)(nil),
	(*Each)(nil),
	(*Len)(nil),
	(*Map)(nil),
	(*Pop)(nil),
	(*Push)(nil),
	(*Reduce)(nil),
	(*Set)(nil),
	(*Slice)(nil),
	(*Splice)(nil),
	// flags:
	(*Case)(nil),
	(*Edge)(nil),
	(*Order)(nil),
	// put:
	(*PutAtEdge)(nil),
	(*PutAtIndex)(nil),
	(*IntoNumList)(nil),
	(*IntoRecList)(nil),
	(*IntoTxtList)(nil),
	//
	(*Sort)(nil),
	(*SortNumbers)(nil),
	(*SortText)(nil),
	(*SortUsing)(nil),
}

func cmdError(op composer.Composer, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}

// can add be inserted into els?
func IsInsertable(ins, els g.Value) (okay bool) {
	return isInsertable(els, ins.Affinity(), ins.Type())
}

// can add be appended to els?
// this is similar to IsInsertable, except that the add can itself be a list.
func IsAppendable(ins, els g.Value) (okay bool) {
	inAff := ins.Affinity()
	if unlist := affine.Element(inAff); len(unlist) > 0 {
		inAff = unlist
	}
	return isInsertable(els, inAff, ins.Type())
}

func isInsertable(els g.Value, haveAff affine.Affinity, haveType string) (okay bool) {
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
	ins, els g.Value
}

func (e insertError) Error() string {
	return errutil.Sprintf("%s of %q isn't insertable into %s of %q",
		e.ins.Affinity(), e.ins.Type(),
		e.els.Affinity(), e.els.Type())
}
