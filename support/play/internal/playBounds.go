package internal

import (
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

type Bounds struct {
	run          rt.Runtime
	focus        string
	bounds       string // name of pattern
	relation     string // parent relation
	lastLocation string // debugging only
}

func MakeDefaultBounds(run rt.Runtime) Bounds {
	return MakeBounds(run, "player", "parser bounds", "whereabouts")
}

// todo: the "patter" refers to bounds for the player
// ideally all of the bounds functions would be in script
// and neither "player" nor "relation" would live here.
func MakeBounds(run rt.Runtime, player, pattern, relation string) Bounds {
	pat, e := run.GetKindByName(pattern)
	if e != nil || !pat.Implements(kindsOf.Pattern.String()) {
		panic(errutil.Sprintf("couldn't find bounds query %q", pattern, e))
	}
	rel, e := run.GetKindByName(relation)
	if e != nil || !rel.Implements(kindsOf.Relation.String()) {
		panic(errutil.Sprintf("couldn't find relation %q", relation, e))
	}
	return Bounds{
		run:      run,
		focus:    player,
		bounds:   pat.Name(),
		relation: rel.Name(),
	}
}

// return the common name of the focal object
// ex. player
func (s *Bounds) GetFocus() string {
	return s.focus
}

// return the id of the focal object
// ex. self
func (s *Bounds) GetFocalObject() (g.Value, error) {
	return s.run.GetField(s.focus, "pawn")
}

// fix: PlayerBounds, PlayerLocale and ObjectBounds might be better delegated to scris.
// it would use the same bits as "locationBounded", only all the bounds requests would use it.
// the scris would switch on the passed string similar to this --
// ( one string and assume that the parser always refers to whatever actor in the global player variable )
// that would probably narrow the dependency on rt -- maybe just to a "call" that could be configured with rt externally.
// ( the "Scope" command could request the named pattern to ensure it exists. )
func (s *Bounds) GetBounds(who, where string) (ret parser.Bounds, err error) {
	run := s.run
	switch who {
	case "", s.focus:
		switch where {
		case "":
			ret, err = s.getLocale()
		case "compass":
			ret = s.compassBounds()
		case s.focus:
			ret = s.selfBounded() // only includes the player's pawn
		default:
			err = errutil.New("unknown bounds", who, where)
		}

	default:
		if _, e := run.GetField(meta.ObjectId, who); e != nil || len(where) > 0 {
			err = errutil.New("unknown bounds", who, where)
		} else {
			ret = s.getObjectBounds(who)
		}
	}
	return
}

// fix: this assumes all objects are emsy
// add containment, whatever...
func (s *Bounds) getObjectBounds(obj string) parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		return
	}
}

func (s *Bounds) getLocale() (ret parser.Bounds, err error) {
	run := s.run
	if pawn, e := s.GetFocalObject(); e != nil {
		err = e
	} else if res, e := run.ReciprocalsOf(pawn.String(), s.relation); e != nil {
		err = e
	} else {
		var where = "nowhere!"
		if res.Len() > 0 {
			v := res.Index(0)
			where = v.String()
		}
		if where != s.lastLocation {
			// log.Println("# GetBounds updated location ", s.player, pawn, where)
			s.lastLocation = where
		}
		ret = s.locationBounded(where) // calls the bounds pattern to return nouns near to the player.
	}
	return
}

// note: play has to attach to some specifics of a particular runtime to work
// ex. childrenOf, transparentOf, etc.
// enc is the enclosureOf the player
func (s *Bounds) locationBounded(enc string) parser.Bounds {
	run := s.run
	return func(cb parser.NounVisitor) (ret bool) {
		if kids, e := run.Call(
			s.bounds,
			affine.TextList,
			[]string{"obj"},
			[]g.Value{g.StringOf(enc)},
		); e != nil && !errors.Is(e, rt.NoResult) {
			log.Println(e)
		} else {
			ret = s.visitStrings(cb, kids)
		}
		return
	}
}

// we have a custom compass bounds because otherwise
// going a direction ( ex. "> north" ) is going to be super slow
func (s *Bounds) compassBounds() (ret parser.Bounds) {
	run := s.run
	return func(cb parser.NounVisitor) (ret bool) {
		if kids, e := run.GetField(meta.ObjectsOfKind, "directions"); e != nil {
			log.Println(e)
		} else {
			ret = s.visitStrings(cb, kids)
		}
		return
	}
}

// return bounds which includes only the player agent and nothing else.
// fix: why this is the agent and not the pawn?
func (s *Bounds) selfBounded() (ret parser.Bounds) {
	return func(cb parser.NounVisitor) (ret bool) {
		return cb(&Noun{s.run, s.focus})
	}
}

func (s *Bounds) visitStrings(cb parser.NounVisitor, kids g.Value) (ret bool) {
	run := s.run
	for _, k := range kids.Strings() {
		if ok := cb(MakeNoun(run, k)); ok {
			ret = ok
			break
		}
	}
	return
}
