package play

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/format"
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

func NewPlaytime(run rt.Runtime, survey Survey, grammar []parser.Scanner) *Playtime {
	g := &parser.AnyOf{Match: grammar}
	return &Playtime{Runtime: run, grammar: g, survey: survey}
}

func (pt *Playtime) Survey() *Survey {
	return &pt.survey
}

type Result struct {
	Action string
	Nouns  []string
}

func (pt *Playtime) RunPlayerAction(name string) (err error) {
	_, err = pt.Call(name, affine.None, nil, []rt.Value{pt.survey.GetFocalObject()})
	return
}

// advance time
func (pt *Playtime) HandleTurn(words string) (ret *Result, err error) {
	var menu format.MenuData // swap and reset
	format.CurrentMenu, menu = menu, format.CurrentMenu
	if len(menu.Action) > 0 {
		ret, err = pt.handleMenus(menu, words)
	} else {
		ret, err = pt.handlePhrases(words)
	}
	return
}

func (pt *Playtime) handleMenus(menu format.MenuData, w string) (ret *Result, err error) {
	str, _ := menu.Match(w)
	nouns := []string{str} // .play adds the player's actor in automatically.
	if e := pt.play(menu.Action, nouns, nil); e != nil {
		err = e
	} else {
		// fix: shouldnt play return result? ( do we even need result? )
		ret = &Result{
			Action: menu.Action,
			Nouns:  nouns,
		}
	}
	return
}

func (pt *Playtime) handlePhrases(words string) (ret *Result, err error) {
	w := pt.Writer()
	switch res, e := pt.scan(words); e.(type) {
	default:
		err = errutil.New("unhandled error", e)

	//"couldnt determine object", a.Nouns)
	case parser.AmbiguousObject:
		fmt.Fprintln(w, e)
	// move to the next state
	// prompt the user, and add whatever the user says into the original input for reparsing
	// insert resolution into input.
	// i, s := e.Depth, append(in, "")
	// copy(s[i+1:], s[i:])
	// s[i] = clarify.NounInstance
	// // println(strings.Join(s, "\\"))
	// err = innerParse(log, pt, match, s, goals)

	// "mismatched word %s != %s at %d", a.Have, a.Want, a.Depth)
	case parser.MismatchedWord:
		fmt.Fprintln(w, "That's not anything i recognize.")

	case parser.MissingObject:
		// in this case, inform guesses at the object to fill.
		fmt.Fprintln(w, e)

	case parser.NoSuchObjects, parser.Overflow, parser.Underflow, parser.UnknownObject:
		fmt.Fprintln(w, e)

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
	if outOfWorld := strings.HasPrefix(act, OutOfWorldPrefix); outOfWorld {
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
			_, err = pt.Runtime.Call(PassTime, affine.None, nil, nil)
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
					_, err = pt.Runtime.Call(PassTime, affine.None, nil, nil)
				}
			}
		}
	}
	return
}

// generic catch all action
func raiseRunAction(run rt.Runtime, actor rt.Value, act string, nouns []string) (okay bool, err error) {
	keys := []string{Actor, Action, FirstNoun, SecondNoun}
	values := []rt.Value{actor, rt.StringOf(act), nounIndex(nouns, 0), nounIndex(nouns, 1)}
	if v, e := run.Call(RunningAnAction, affine.None, keys, values); e != nil {
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
