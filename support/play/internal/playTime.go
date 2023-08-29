package internal

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/parser"
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
	bounds   string
}

func NewPlaytime(run *qna.Runner) *Playtime {
	return NewCustomPlaytime(run, "player", "whereabouts", lang.Normalize("parser bounds"))
}

// the named relation should yield a single object for the named player.
// the bounds pattern should return the objects in that player's local area.
func NewCustomPlaytime(run *qna.Runner, player, relation, bounds string) *Playtime {
	if _, e := run.GetKindByName(bounds); e != nil {
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
// future: to differentiate b/t system actions and "timed" actions,
// consider using naming convention: ex. @save (mud style), or #save
func (pt *Playtime) Play(name, player string, args []string) (err error) {
	// temp patch to new:
	// should instead raise a parsing event with the nouns and the action name
	// ( possibly -- probably send in the player since it would be needed for bounds still )
	// player should be a variable not a pawn; although handling the difference here helps.
	if actor, e := pt.GetField(player, "pawn"); e != nil {
		err = e
	} else {
		// insert the player in front of the other args.
		vs := make([]g.Value, len(args)+1)
		vs[0] = actor
		for i, n := range args {
			vs[i+1] = g.StringOf(n)
		}
		_, err = pt.Call(name, affine.None, nil, vs)
	}
	return
}

func (pt *Playtime) IsPlural(word string) bool {
	pl := pt.SingularOf(word)
	return len(pl) > 0 && pl != word
}

var lastLocation string // debugging only

// fix: PlayerBounds, PlayerLocale and ObjectBounds might be better delegated to script.
// it would use the same bits as "locationBounded", only all the bounds requests would use it.
// the script would switch on the passed string similar to this --
// ( one string and assume that the parser always refers to whatever actor in the global player variable )
// that would probably narrow the dependency on rt -- maybe just to a "call" that could be configured with rt externally.
// ( the "Scope" command could request the named pattern to ensure it exists. )
func (pt *Playtime) GetPlayerBounds(where string) (ret parser.Bounds, err error) {
	switch where {
	case "":
		ret, err = pt.GetPlayerLocale()
	case "player":
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
			// log.Println("# GetPlayerBounds updated location ", pt.player, pawn, where)
			lastLocation = where
		}
		ret = pt.locationBounded(where) // calls the bounds pattern to return nouns near to the player.
	}
	return
}

// fix: assumes all objects are empty
// add containment, whatever...
func (pt *Playtime) GetObjectBounds(obj string) (ret parser.Bounds, err error) {
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
