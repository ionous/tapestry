package play

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

// Playtime - adapts the qna.Runner rt.Runtime to the parser
// this is VERY rudimentary.
type Playtime struct {
	rt.Runtime
	grammar parser.Scanner
	survey  Survey
}

func NewPlaytime(run rt.Runtime, survey Survey, grammar parser.Scanner) *Playtime {
	return &Playtime{Runtime: run, grammar: grammar, survey: survey}
}

type Result struct {
	Action string
	Nouns  []string
}

// advance time
func (pt *Playtime) Step(words string) (ret *Result, err error) {
	switch res, e := pt.scan(words); e.(type) {
	default:
		err = errutil.New("unhandled error", e)

	//"couldnt determine object", a.Nouns)
	// case parser.AmbiguousObject:
	// move to the next state
	// prompt the user, and add whatever the user says into the original input for reparsing
	// insert resolution into input.
	// i, s := e.Depth, append(in, "")
	// copy(s[i+1:], s[i:])
	// s[i] = clarify.NounInstance
	// // println(strings.Join(s, "\\"))
	// err = innerParse(log, pt, match, s, goals)

	// "mismatched word %s != %s at %d", a.Have, a.Want, a.Depth)
	// case parser.MismatchedWord:

	// "missing an object at %d"
	// case parser.MissingObject:
	// in this case, inform guesses at the object to fill.

	// "you cant see any such things"
	// case parser.NoSuchObjects:

	// "too many words"
	// case parser.Overflow:

	// "too few words"
	// case parser.Underflow:

	// "you can't see any such thing"
	case parser.UnknownObject:
		fmt.Println(e)

	case nil:
		switch res := res.(type) {
		// usually, we get a result list
		// the last element of which is an action
		case *parser.ResultList:
			if last, ok := res.Last(); !ok {
				err = errutil.New("result list was empty")
			} else if act, ok := last.(parser.ResolvedAction); !ok {
				err = errutil.Fmt("expected resolved action %T", last)
			} else {
				// multi-actions are probably repeats or something
				// or maybe get passed lists of objects hrmm.
				// send these nouns to the runtime
				nouns := res.Objects()
				if e := pt.play(act.Name, nouns, act.Args); e != nil {
					err = errutil.Fmt("%w for %v", e, res)
				} else {
					ret = &Result{
						Action: act.Name,
						Nouns:  nouns,
					}
				}
			}

		// - Action terminates a matcher sequence, resolving to the named action.
		// case parser.ResolvedAction:
		//- Multi matches one or more objects. ( the words for "all, etc." are hard-coded )
		// case parser.ResolvedMulti:

		// - Noun matches one object held by the context.
		// case parser.ResolvedNoun:

		// - Word matches one word.
		// case parser.ResolvedWords:
		default:
			err = errutil.New("unhandled results", res)
		} // end res
	} // end err
	return
} // end func

func (pt *Playtime) scan(words string) (ret parser.Result, err error) {
	if bounds, e := pt.survey.GetBounds("", ""); e != nil {
		err = e
	} else {
		ctx := (*parserContext)(pt)
		cursor := parser.Cursor{Words: strings.Fields(words)}
		ret, err = pt.grammar.Scan(ctx, bounds, cursor)
	}
	return
}

// execute a command command
// future: to differentiate b/t system actions and "timed" actions,
// consider using naming convention: ex. @save (mud style), or #save
func (pt *Playtime) play(act string, nouns []string, args []assign.Arg) (err error) {
	// fix: raise a parsing event with the nouns and the action name
	if focus := pt.survey.GetFocalObject(); focus == nil {
		err = errutil.New("couldnt get focal object")
	} else if ks, vs, e := assign.ExpandArgs(pt, args); e != nil {
		err = e
	} else {
		// the actor ( and any nouns ) need to precede the "keyed" fields.
		els := make([]g.Value, 1, 1+len(nouns)+len(vs))
		els[0] = focus // presumably the player's actor
		for _, n := range nouns {
			els = append(els, g.StringOf(n))
		}
		els = append(els, vs...)
		if _, e := pt.Runtime.Call(act, affine.None, ks, els); e != nil {
			err = e
		} else if !strings.HasPrefix(act, "requesting ") {
			// meta actions are defined as those with start with the string "requesting ..."
			if _, e := pt.Runtime.Call("pass time", affine.None, nil, els[:1]); e != nil {
				err = e
			}
		}
	}
	return
}
