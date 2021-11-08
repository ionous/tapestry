package internal

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/qna"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// Playtime - adapts the qna.Runner rt.Runtime to the parser
// this is VERY rudimentary.
type Playtime struct {
	*qna.Runner
	hasName  *sql.Stmt
	player   ident.Id
	location string
}

func NewPlaytime(db *sql.DB, player, startWhere string) *Playtime {
	var ps tables.Prep
	return &Playtime{
		Runner:   qna.NewRuntime(db, iffy.AllSignatures),
		location: startWhere,
		player:   ident.IdOf(player),
		hasName: ps.Prep(db,
			`select 1 
				from mdl_name
				where noun= ?1 and name = ?2`),
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
	switch e := pt.hasName.QueryRow(noun, name).Scan(&ret); e {
	default:
		panic(e)
	case nil, sql.ErrNoRows:
		// use scanned result
	}
	return
}
