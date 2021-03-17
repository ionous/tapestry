package internal

import (
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type Parser struct {
	pt   *Playtime
	gram parser.Scanner
	// state    int or something
}

func NewParser(pt *Playtime, gram parser.Scanner) *Parser {
	if gram == nil {
		gram = Grammar
	}
	return &Parser{pt, gram}
}

type Result struct {
	Action string
	Nouns  []string
}

func (p *Parser) Step(words string) (ret *Result, err error) {
	pt := p.pt
	bounds := pt.GetNamedBounds(pt.location)
	cursor := parser.Cursor{Words: strings.Fields(words)}
	switch res, e := p.gram.Scan(pt, bounds, cursor); e.(type) {
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
	// case parser.UnknownObject:

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
				objs := res.Objects()
				var args argList
				out := new(Result)
				out.Action = act.Name
				if cnt := len(objs); true {
					args = make([]rt.Arg, 0, cnt+1)
					args.add(
						&core.FromText{
							&core.IdOf{
								&core.Text{Text: "player"},
							}})
					//
					out.Nouns = make([]string, cnt)
					for i, obj := range objs {
						n := obj.String()
						out.Nouns[i] = n
						args.add(
							&core.FromValue{g.StringOf(n)},
						)
					}
				}
				// send these nouns to the runtime
				if _, e := pt.Call(act.Name, "", args); e != nil {
					err = errutil.New(e, "for", res)
				} else {
					ret = out
				}
			}

		// - Action terminates a matcher sequence, resolving to the named action.
		// case parser.ResolvedAction:
		//- Multi matches one or more objects. ( the words for "all, etc." are hard-coded )
		// case parser.ResolvedMulti:

		// - Noun matches one object held by the context.
		// case parser.ResolvedNoun:

		// - Word matches one word.
		// case parser.ResolvedWord:
		default:
			err = errutil.New("unhandled results", res)
		} // end res
	} // end err
	return
} // end func

type argList []rt.Arg

func (l *argList) add(a rt.Assignment) {
	(*l) = append((*l), rt.Arg{
		Name: "$" + strconv.Itoa(len((*l))+1),
		From: a,
	})
}
