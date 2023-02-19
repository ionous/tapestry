package internal

import (
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/parser/ident"
	"git.sr.ht/~ionous/tapestry/qna"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

// Playtime - adapts the qna.Runner rt.Runtime to the parser
// this is VERY rudimentary.
type Playtime struct {
	*qna.Runner
	player   string
	relation string
	bounds   *g.Kind
}

func NewPlaytime(run *qna.Runner) *Playtime {
	return NewCustomPlaytime(run, "player", "whereabouts", "parser_bounds")
}

// the named relation should yield a single object for the named player.
// the bounds pattern should return the objects in that player's local area.
func NewCustomPlaytime(run *qna.Runner, player, relation, bounds string) *Playtime {
	if bounds, e := run.GetKindByName(bounds); e != nil {
		panic(e)
	} else {
		return &Playtime{
			Runner:   run,
			player:   player,
			relation: relation,
			bounds:   bounds,
		}
	}
}

// step the world by running some command
func (pt *Playtime) Play(name, player string, args []string) (err error) {
	// future: to differentiate b/t system actions and "timed" actions,
	// consider using naming convention: ex. #save.
	if k, e := pt.GetKindByName(name); e != nil {
		err = e
	} else if max, min := k.NumField(), len(args)+1; max < min {
		err = errutil.New("not enough fields", min, max)
	} else {
		rec := k.NewRecord()
		for i, a := range append([]string{player}, args...) {
			f := k.Field(i)
			// fix conversion of numbers
			if aff := f.Affinity; aff != affine.Text {
				err = errutil.New("field", i, "expected text, have", aff)
				break
			} else if e := rec.SetIndexedField(i, g.StringOf(a)); e != nil {
				err = errutil.New("field", i, e)
				break
			}
		}
		if err == nil {
			if _, e := pt.Call(rec, affine.None); e != nil {
				err = e
			}
		}
	}
	return
}

func (pt *Playtime) IsPlural(word string) bool {
	pl := pt.SingularOf(word)
	return len(pl) > 0 && pl != word
}

var lastLocation string // debugging only

func (pt *Playtime) GetPlayerBounds(where string) (ret parser.Bounds, err error) {
	switch where {
	case "":
		ret, err = pt.GetPlayerLocale()
	case "self":
		ret, err = pt.selfBounded() // only includes the player's pawn
	default:
		err = errutil.New("unknown player bounds", where)
	}
	return
}

func (pt *Playtime) GetPlayerLocale() (ret parser.Bounds, err error) {
	if pawn, e := pt.getPawn(); e != nil {
		err = e
	} else if res, e := pt.ReciprocalsOf(pawn, pt.relation); e != nil {
		err = e
	} else {
		var where = "nowhere!"
		if res.Len() > 0 {
			v := res.Index(0)
			where = v.String()
		}
		if where != lastLocation {
			log.Println("GetPlayerBounds", pt.player, pawn, where)
			lastLocation = where
		}
		ret = pt.locationBounded(where) // calls the bounds pattern to return nouns near to the player.
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
