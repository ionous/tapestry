package internal

import (
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/parser/ident"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Playtime - adapts the qna.Runner rt.Runtime to the parser
// this is VERY rudimentary.
type Playtime struct {
	*qna.Runner
	player   ident.Id
	location string
}

func NewPlaytime(run *qna.Runner, player, startWhere string) *Playtime {
	return &Playtime{
		Runner:   run,
		location: startWhere,
		player:   ident.IdOf(player),
	}
}

func (pt *Playtime) Play(name string, args []rt.Arg) (err error) {
	// future: to differentiate b/t system actions and "timed" actions,
	// consider using naming convention: ex. #save.
	if _, e := pt.Call(name, "", args); e != nil {
		err = e
	}
	return
}

func (pt *Playtime) IsPlural(word string) bool {
	pl := pt.SingularOf(word)
	return len(pl) > 0 && pl != word
}

func (pt *Playtime) GetPlayerBounds(where string) (ret parser.Bounds, err error) {
	switch where {
	case "":
		ret = pt.LocationBounded(pt.location)
	case "self":
		ret = pt.SelfBounded()
	default:
		err = errutil.New("unknown player bounds", where)
	}
	return
}

// fix: assumes all objects are empty
// add containment, whatever...
func (pt *Playtime) GetObjectBounds(obj ident.Id) (ret parser.Bounds, err error) {
	ret = func(cb parser.NounVisitor) (ret bool) {
		return
	}
	return
}

func (pt *Playtime) HasName(noun, name string) (ret bool) {
	if ok, e := pt.NounIsNamed(noun, name); e != nil {
		panic(e)
	} else {
		ret = ok
	}
	return
}
