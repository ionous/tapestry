package play

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
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
	// state    int or something
}

func NewPlaytime(run rt.Runtime, grammar parser.Scanner) *Playtime {
	survey := MakeDefaultSurveyor(run)
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
				if e := pt.play(act.Name, nouns); e != nil {
					err = errutil.New(e, "for", res)
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
func (pt *Playtime) play(act string, args []string) (err error) {
	// temp patch to new:
	// should instead raise a parsing event with the nouns and the action name
	// ( possibly -- probably send in the player since it would be needed for bounds still )
	if actor, e := pt.survey.GetFocalObject(); e != nil {
		err = e
	} else {
		// insert the player in front of the other args.
		vs := make([]g.Value, len(args)+1)
		vs[0] = actor
		for i, n := range args {
			vs[i+1] = g.StringOf(n)
		}
		_, err = pt.Runtime.Call(act, affine.None, nil, vs)
	}
	return
}
