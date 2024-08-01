package play

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/literal"
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

func (pt *Playtime) RunPlayerAction(name string) (err error) {
	_, err = pt.Call(name, affine.None, nil, []rt.Value{pt.survey.GetFocalObject()})
	return
}

// advance time
func (pt *Playtime) HandleTurn(words string) (okay bool, err error) {
	var menu format.MenuData // swap and reset
	format.CurrentMenu, menu = menu, format.CurrentMenu
	if len(menu.Action) > 0 {
		okay, err = pt.handleMenus(menu, words)
	} else {
		okay, err = pt.handlePhrases(words)
	}
	return
}

func (pt *Playtime) handleMenus(menu format.MenuData, w string) (okay bool, err error) {
	str, _ := menu.Match(w)
	if e := pt.play(menu.Action, nil, []call.Arg{{
		Value: &call.FromText{Value: literal.T(str)}},
	}); e != nil || len(w) == 0 {
		err = e // loop on invalid choices
		format.CurrentMenu = menu
	} else {
		okay = true
	}
	return
}

func (pt *Playtime) handlePhrases(words string) (okay bool, err error) {
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
					okay = true
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
	// 1. out of world requests:
	if outOfWorld := strings.HasPrefix(act, OutOfWorldPrefix); outOfWorld {
		if len(nouns) != 0 {
			err = errutil.New("out of world actions don't expect any nouns")
		} else if ks, vs, e := call.ExpandArgs(pt, args); e != nil {
			err = e
		} else {
			_, err = pt.Runtime.Call(act, affine.None, ks, vs)
		}
	} else {
		// player action:
		if focus := pt.survey.GetFocalObject(); focus == nil {
			err = errutil.New("couldn't get focal object")
		} else if ok, e := raiseRunAction(pt, focus, act, nouns); e != nil {
			err = e
		} else {
			// "running an action" returned true
			// permitting us to call the requested parser action:
			if ok {
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
					_, err = pt.Runtime.Call(act, affine.None, ks, els)
				}
			}
			// pass time.
			if err == nil {
				_, err = pt.Runtime.Call(PassTime, affine.None, nil, nil)
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
