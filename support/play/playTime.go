package play

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/rt"
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

func (p *Playtime) Survey() *Survey {
	return &p.survey
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
		fmt.Println() // command break

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

// execute a command
func (pt *Playtime) play(act string, nouns []string, args []call.Arg) (err error) {
	if outOfWorld := strings.HasPrefix(act, "request "); outOfWorld {
		if len(nouns) != 0 {
			// fix: check at weave?
			err = errutil.New("out of world actions don't expect any nouns")
		} else if ks, vs, e := call.ExpandArgs(pt, args); e != nil {
			err = e
		} else {
			_, err = pt.Runtime.Call(act, affine.None, ks, vs)
		}
	} else {
		// fix: raise a parsing event with the nouns and the action name
		if focus := pt.survey.GetFocalObject(); focus == nil {
			err = errutil.New("couldnt get focal object")
		} else if ok, e := raiseRunAction(pt, focus, act, nouns); e != nil {
			err = e
		} else if !ok {
			_, err = pt.Runtime.Call("pass time", affine.None, nil, nil)
		} else {
			if ks, vs, e := call.ExpandArgs(pt, args); e != nil {
				err = e
			} else {
				// the actor ( and any nouns ) need to precede the "keyed" fields.
				els := make([]rt.Value, 1, 1+len(nouns)+len(vs))
				els[0] = focus // presumably the player's actor
				for _, n := range nouns {
					els = append(els, rt.StringOf(n))
				}
				els = append(els, vs...)
				if _, e := pt.Runtime.Call(act, affine.None, ks, els); e != nil {
					err = e
				} else {
					_, err = pt.Runtime.Call("pass time", affine.None, nil, nil)
				}
			}
		}
	}
	return
}

// generic catch all action
func raiseRunAction(run rt.Runtime, actor rt.Value, act string, nouns []string) (okay bool, err error) {
	keys := []string{"actor", "action", "first noun", "second noun"}
	values := []rt.Value{actor, rt.StringOf(act), nounIndex(nouns, 0), nounIndex(nouns, 1)}
	if v, e := run.Call("running an action", affine.None, keys, values); e != nil {
		err = e
	} else {
		okay = v.Bool()
	}
	return
}

func nounIndex(nouns []string, i int) (ret rt.Value) {
	if i < len(nouns) {
		ret = rt.StringOf(nouns[i])
	} else {
		ret = rt.Nothing
	}
	return
}
