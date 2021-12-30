package parser_test

import (
	"strings"

	. "git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/parser/ident"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func actor() Scanner {
	return noun(&HasClass{"actors"})
}

func anyOf(s ...Scanner) (ret Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &AnyOf{s}
	}
	return
}

func allOf(s ...Scanner) (ret Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &AllOf{s}
	}
	return
}

func noun(f ...Filter) Scanner {
	return &Noun{f}
}
func nouns(f ...Filter) Scanner {
	return &Multi{f}
}

// changes the bounds of its first scanner in response to the results of its last scanner.
// func retarget(s ...Scanner) Scanner {
// 	return &Target{s}
// }

// swaps the first detected noun with the second detected noun
func reverse(s ...Scanner) Scanner {
	return &Reverse{s}
}

// match any word/phrase referring to the player:
// forces the focus to just the player object,
// attempts to match the specified noun to that player object,
// absorbs all the words that were used to match the player.
func self() Scanner {
	return &Eat{&Focus{"self", &Noun{}}}
}

// note: we use things to exclude directions
func thing() Scanner {
	return noun(&HasClass{"things"})
}

func things() Scanner {
	return nouns(&HasClass{"things"})
}

func words(s ...string) (ret Scanner) {
	return Words(s)
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

var showGrammar = allOf(words("show"), anyOf(
	allOf(noun(), words("to"), self(), &Action{"Examine"}),
	allOf(noun(), words("to"), actor(), &Action{"Show"}),
	allOf(reverse(actor(), noun()), &Action{"Show"}),
))

func makeObject(name string, kinds ...string) *MyObject {
	names := strings.Fields(name)
	if kinds[0] != "actors" {
		kinds = append(kinds, "things")
	}
	id := ident.IdOf(strings.Join(names, "-"))
	return &MyObject{Id: id, Names: names, Classes: kinds}
}

// MyBounds implements parser.Context
var ctx = func() (ret MyBounds) {
	ret = MyBounds{
		makeObject("something", "things"),
		makeObject("red apple", "apples"),
		makeObject("crab apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
		makeObject("torch", "devices"),
		makeObject("bob", "actors"),
	}
	return append(ret, Directions...)
}()

var myself = makeObject("self", "actors")

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
