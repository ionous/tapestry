package internal

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// Playtime - adapts the qna.Runner rt.Runtime to the parser
// this is VERY rudimentary.
type Playtime struct {
	rt.Runtime
	Bounds
}

func NewPlaytime(run rt.Runtime) *Playtime {
	bounds := MakeDefaultBounds(run)
	return NewCustomPlaytime(run, bounds)
}

// the named relation should yield a single object for the named player.
// the bounds pattern should return the objects in that player's local area.
func NewCustomPlaytime(run rt.Runtime, bounds Bounds) *Playtime {
	return &Playtime{
		Runtime: run,
		Bounds:  bounds,
	}
}

// step the world by running some command
// future: to differentiate b/t system actions and "timed" actions,
// consider using naming convention: ex. @save (mud style), or #save
func (pt *Playtime) Play(act string, args []string) (err error) {
	// temp patch to new:
	// should instead raise a parsing event with the nouns and the action name
	// ( possibly -- probably send in the player since it would be needed for bounds still )
	if actor, e := pt.GetFocalObject(); e != nil {
		err = e
	} else {
		// insert the player in front of the other args.
		vs := make([]g.Value, len(args)+1)
		vs[0] = actor
		for i, n := range args {
			vs[i+1] = g.StringOf(n)
		}
		_, err = pt.Call(act, affine.None, nil, vs)
	}
	return
}

func (pt *Playtime) IsPlural(word string) bool {
	pl := pt.SingularOf(word)
	return len(pl) > 0 && pl != word
}
