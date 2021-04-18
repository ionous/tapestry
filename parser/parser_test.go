package parser_test

import (
	"strings"

	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/parser"
	. "git.sr.ht/~ionous/iffy/parser"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func anyOf(s ...parser.Scanner) (ret parser.Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &parser.AnyOf{s}
	}
	return
}

func allOf(s ...parser.Scanner) (ret parser.Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &parser.AllOf{s}
	}
	return
}

func noun(f ...Filter) parser.Scanner {
	return &parser.Noun{f}
}
func nouns(f ...Filter) parser.Scanner {
	return &parser.Multi{f}
}

// changes the bounds of its first scanner in response to the results of its last scanner.
func retarget(s ...parser.Scanner) parser.Scanner {
	return &parser.Target{s}
}

// note: we use things to exclude directions
func thing() parser.Scanner {
	return noun(&parser.HasClass{"things"})
}

func things() parser.Scanner {
	return nouns(&parser.HasClass{"things"})
}

func words(s ...string) (ret parser.Scanner) {
	return parser.Words(s)
}

var lookGrammar = allOf(words("look", "l"), anyOf(
	allOf(&Action{"Look"}),
	allOf(words("at"), noun(), &Action{"Examine"}),
	// before "look inside", since inside is also direction.
	allOf(noun(&HasClass{"directions"}), &Action{"Examine"}),
	allOf(words("to"), noun(&HasClass{"directions"}), &Action{"Examine"}),
	allOf(words("inside", "in", "into", "through", "on"), noun(), &Action{"Search"}),
	allOf(words("under"), noun(), &Action{"LookUnder"}),
))

var pickGrammar = allOf(words("pick"), anyOf(
	allOf(words("up"), things(), &Action{"Take"}),
	allOf(things(), words("up"), &Action{"Take"}),
))

func makeObject(s ...string) *MyObject {
	name, s := s[0], s[1:]
	names := strings.Fields(name)
	s = append(s, "things")
	id := ident.IdOf(strings.Join(names, "-"))
	return &MyObject{Id: id, Names: names, Classes: s}
}

var ctx = func() (ret MyBounds) {
	ret = MyBounds{
		makeObject("something"),
		makeObject("red apple", "apples"),
		makeObject("crab apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
		makeObject("torch", "devices"),
	}
	return append(ret, Directions...)
}()

// StringIds - convert a list of strings to ids
func StringIds(strs []string) (ret []ident.Id) {
	for _, str := range strs {
		ret = append(ret, ident.IdOf(str))
	}
	return
}

// Log helper - matches testing.T
type Log interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func parse(log Log, ctx Context, match Scanner, phrases []string, goals ...Goal) (err error) {
	for _, in := range phrases {
		fields := strings.Fields(in)
		if e := innerParse(log, ctx, match, fields, goals); e != nil {
			err = errutil.Fmt("%v for '%s'", e, in)
			break
		}
	}
	return
}

// test a string of words against a grammar to see if they match the desired goals
func innerParse(log Log, ctx Context, match Scanner, in []string, goals []Goal) (err error) {
	if len(goals) == 0 {
		err = errutil.New("expected some goals")
	} else {
		goal, goals := goals[0], goals[1:]
		if bounds, e := ctx.GetPlayerBounds(""); e != nil {
			err = e
		} else if res, e := match.Scan(ctx, bounds, Cursor{Words: in}); e != nil {
			// on error:
			switch want := goal.(type) {
			case *ErrorGoal:
				if e.Error() != want.Error {
					err = errutil.Fmt("mismatched error want:'%s' got:'%s'", want, e)
				} else {
					log.Log("matched error", []error{e})
				}
			case *ClarifyGoal:
				clarify := want
				switch e := e.(type) {
				case MissingObject:
					extend := append(in, clarify.NounInstance)
					err = innerParse(log, ctx, match, extend, goals)
				case AmbiguousObject:
					// println(strings.Join(in, "/"))
					// insert resolution into input.
					i, s := e.Depth, append(in, "")
					copy(s[i+1:], s[i:])
					s[i] = clarify.NounInstance
					// println(strings.Join(s, "\\"))
					err = innerParse(log, ctx, match, s, goals)
				default:
					err = errutil.Fmt("clarification not implemented for %T", e)
				}
			default:
				err = errutil.New("unexpected failure:", e)
			}
		} else if goal == nil {
			err = errutil.New("unexpected success")
		} else if want, ok := goal.(*ActionGoal); !ok {
			err = errutil.Fmt("unexpected goal %s %T for result %v", in, goal, pretty.Sprint(res))
		} else if results, ok := res.(*ResultList); !ok {
			err = errutil.Fmt("expected result list %T", res)
		} else if last, ok := results.Last(); !ok {
			err = errutil.New("result list was empty")
		} else if act, ok := last.(ResolvedAction); !ok {
			err = errutil.Fmt("expected resolved action %T", last)
		} else if !strings.EqualFold(act.Name, want.Action) {
			err = errutil.New("expected action", act, "got", want.Action)
		} else if want, have := PrettyIds(want.Objects()), PrettyIds(results.Objects()); want != have {
			err = errutil.Fmt("expected nouns %q got %q", want, have)
		} else {
			log.Logf("matched %v", in)
		}
	}
	return
}
