package play

import (
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

// the survey creates bounds: areas of the world containing sets of nouns.
type Survey struct {
	run          rt.Runtime
	focus        string
	bounds       string // name of pattern
	relation     string // parent relation
	lastLocation string // debugging only
}

func MakeDefaultSurveyor(run rt.Runtime) Survey {
	return MakeSurveyor(run, "player", "survey bounds", "whereabouts")
}

// todo: the "patter" refers to bounds for the player
// ideally all of the bounds functions would be in script
// and neither "player" nor "relation" would live here.
func MakeSurveyor(run rt.Runtime, player, pattern, relation string) Survey {
	// if we validate upfront we have to activate the domain(s) before creating the surveyor
	// pat, e := run.GetKindByName(pattern)
	// if e != nil || !pat.Implements(kindsOf.Pattern.String()) {
	// 	panic(errutil.Sprintf("couldn't find bounds query %s", e))
	// }
	// rel, e := run.GetKindByName(relation)
	// if e != nil || !rel.Implements(kindsOf.Relation.String()) {
	// 	panic(errutil.Sprintf("couldn't find relation %s", e))
	// }
	return Survey{
		run:      run,
		focus:    player,
		bounds:   lang.Normalize(pattern),
		relation: lang.Normalize(relation),
	}
}

// return the common name of the focal object
// ex. player
func (s *Survey) GetFocus() string {
	return s.focus
}

// return the id of the focal object; returns nil on error; ex. self
func (s *Survey) GetFocalObject() (ret g.Value) {
	ret, _ = s.run.GetField(s.focus, "pawn")
	return
}

// fix: PlayerBounds, PlayerLocale and ObjectBounds might be better delegated to scris.
// it would use the same bits as "locationBounded", only all the bounds requests would use it.
// the scris would switch on the passed string similar to this --
// ( one string and assume that the parser always refers to whatever actor in the global player variable )
// that would probably narrow the dependency on rt -- maybe just to a "call" that could be configured with rt externally.
// ( the "Scope" command could request the named pattern to ensure it exists. )
func (s *Survey) GetBounds(who, where string) (ret parser.Bounds, err error) {
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
func (s *Survey) getObjectBounds(obj string) parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		return
	}
}

func (s *Survey) getLocale() (ret parser.Bounds, err error) {
	run := s.run
	if actor := s.GetFocalObject(); actor == nil {
		err = errutil.New("couldnt get focal object")
	} else if res, e := run.ReciprocalsOf(actor.String(), s.relation); e != nil {
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
func (s *Survey) locationBounded(enc string) parser.Bounds {
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
func (s *Survey) compassBounds() (ret parser.Bounds) {
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
func (s *Survey) selfBounded() (ret parser.Bounds) {
	return func(cb parser.NounVisitor) (ret bool) {
		return cb(&Noun{s.run, s.focus})
	}
}

func (s *Survey) visitStrings(cb parser.NounVisitor, kids g.Value) (ret bool) {
	run := s.run
	for _, k := range kids.Strings() {
		if ok := cb(MakeNoun(run, k)); ok {
			ret = ok
			break
		}
	}
	return
}
