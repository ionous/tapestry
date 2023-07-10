package parser

import (
	"github.com/ionous/errutil"
)

// Focus - a Scanner which changes the bounds for subsequent scanners.
// For instance, searching only though held objects.
type Focus struct {
	Where string
	What  Scanner
}

//
func (a *Focus) Scan(ctx Context, _ Bounds, cs Cursor) (ret Result, err error) {
	if bounds, e := ctx.GetPlayerBounds(a.Where); e != nil {
		err = e
	} else {
		ret, err = a.What.Scan(ctx, bounds, cs)
	}
	return
}

// Target changes the bounds of its first scanner in response to the results of its last scanner.
// Generally, this means that the last scanner should be Noun{}.
type Target struct {
	Match []Scanner
}

//
func (a *Target) Scan(ctx Context, bounds Bounds, start Cursor) (ret Result, err error) {
	first, rest := a.Match[0], a.Match[1:]
	errorDepth := -1
	err = Underflow{Depth(start.Pos)} // not completely sure what a good default error is here...
	// scan ahead for matches and determine how many words might match this target.
	for cs := start; len(cs.CurrentWord()) > 0; cs = cs.Skip(1) {
		if rl, e := scanAllOf(ctx, bounds, cs, rest); e != nil {
			// like anyOf, we track the "deepest" / latest error.
			if d := DepthOf(e); d >= errorDepth {
				err, errorDepth = e, d
			}
			continue // keep looking for success
		} else if last, ok := rl.Last(); !ok {
			err = errutil.New("target not found")
		} else if obj, ok := last.(ResolvedNoun); !ok {
			err = errutil.Fmt("expected an object, got %T", last)
			break
		} else if bounds, e := ctx.GetObjectBounds(obj.NounInstance.String()); e != nil {
			err = e
			break
		} else {
			// snip down our input to just the preface which matches this target
			// ( to avoid errors of "too many words" )
			words := start.Words[:cs.Pos]
			sub := Cursor{start.Pos, words}
			if r, e := first.Scan(ctx, bounds, sub); e != nil {
				err = e
				break
			} else {
				rl.addResult(r)
				ret, err = rl, nil
				break
			}
		}
	}
	return
}
