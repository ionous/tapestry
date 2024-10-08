package parser_test

import (
	"strings"
	"testing"

	. "git.sr.ht/~ionous/tapestry/parser"
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
	"github.com/ionous/sliceOf"
)

// MyObject - along with MyNoun this provides an example of mapping some application defined object to a parser.NounInstance.
type MyObject struct {
	Id         string
	Names      []string
	Classes    []string
	Attributes []string
}

func (m *MyObject) String() string {
	return m.Id
}

type MyBounds []*MyObject

// some of the simpler tests have objects ids which are a single rune long.
func (m MyBounds) Get(r rune) NounInstance {
	return MyNoun{m[r-'a']}
}

// generate a bunch of nouns from single runes for certain tests.
func (m MyBounds) Many(rs ...rune) (ret []NounInstance) {
	for _, r := range rs {
		ret = append(ret, m.Get(r))
	}
	return
}

func (m MyBounds) GetBounds(who, where string) (ret Bounds, err error) {
	switch who {
	default:
		ret = m.SearchBounds
	case "":
		switch where {
		case "":
			ret = m.SearchBounds
		case "player":
			ret = m.PlayerBounds
		default:
			err = errutil.New("unknown bounds", who, where)
		}
	}
	return
}

func (m MyBounds) IsPlural(word string) bool {
	return word != inflect.Singularize(word)
}

func (m MyBounds) PlayerBounds(v NounVisitor) (ret bool) {
	n := MyNoun{MyObject: myself}
	return v(n)
}

func (m MyBounds) SearchBounds(v NounVisitor) (ret bool) {
	for _, k := range m {
		if v(MyNoun{k}) {
			ret = true
			break
		}
	}
	return
}

// MyNoun implements NounInstance for MyObject
type MyNoun struct {
	*MyObject
}

func (adapt MyNoun) Id() string {
	return adapt.MyObject.Id
}

func (adapt MyNoun) HasName(name string) bool {
	return MatchAny(name, adapt.Names)
}

func (adapt MyNoun) HasClass(cls string) bool {
	return MatchAny(cls, adapt.Classes)
}

func (adapt MyNoun) HasPlural(plural string) bool {
	// we'll use classes as plurals for tests --
	// its possible that might be different for the runtime
	// ex. might check plural / printed names
	return MatchAny(plural, adapt.Classes)
}

func (adapt MyNoun) HasAttribute(attr string) bool {
	return MatchAny(attr, adapt.Attributes)
}

func MatchAny(n string, l []string) (okay bool) {
	for _, s := range l {
		if strings.EqualFold(n, s) {
			okay = true
			break
		}
	}
	return
}

func TestBounds(t *testing.T) {
	ctx := MyBounds{
		&MyObject{Id: "a", Names: sliceOf.String("unique")},
		//
		&MyObject{Id: "b", Names: strings.Fields("exact")},
		&MyObject{Id: "c", Names: strings.Fields("exact match")},
		//
		&MyObject{Id: "d", Names: strings.Fields("inexact match")},
		&MyObject{Id: "e", Names: strings.Fields("inexact conflict")},
		//
		&MyObject{Id: "f",
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
			Classes:    strings.Fields("class"),
		},
		&MyObject{Id: "g",
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
		},
		&MyObject{Id: "h",
			Names:   strings.Fields("filter"),
			Classes: strings.Fields("class"),
		},
	}
	if res, e := matching(ctx, "unique"); e != nil {
		t.Fatal("error", e)
	} else if obj, ok := res.(ResolvedNoun); !ok {
		t.Fatalf("%T", res)
	} else if obj.NounInstance != ctx.Get('a') {
		t.Fatal("mismatched", obj.NounInstance)
	} else if got, want := strings.Join(obj.Words, ","), "unique"; got != want {
		t.Fatal(got)
	}

	if res, e := matching(ctx, "exact match"); e != nil {
		t.Fatal("error", e)
	} else if obj, ok := res.(ResolvedNoun); !ok {
		t.Fatalf("%T", res)
	} else if obj.NounInstance != ctx.Get('c') {
		t.Fatal("mismatched", obj.NounInstance)
	} else if got, want := strings.Join(obj.Words, ","), "exact,match"; got != want {
		t.Fatal(got)
	}

	if res, e := matchingFilter(ctx, "filter", "attr", "class"); e != nil {
		t.Fatal("error", e)
	} else if obj, ok := res.(ResolvedNoun); !ok {
		t.Fatalf("%T", res)
	} else if obj.NounInstance != ctx.Get('f') {
		t.Fatal("mismatched", obj.NounInstance)
	} else if got, want := strings.Join(obj.Words, ","), "filter"; got != want {
		t.Fatal(got)
	}

	if res, e := matching(ctx, "inexact"); e == nil || res != nil {
		t.Fatal("expected error", e, res)
	} else if got, want := e.Error(), (AmbiguousObject{
		Nouns: ctx.Many('d', 'e'),
		Depth: 1,
	}).Error(); got != want {
		t.Fatal(got)
	}

	if res, e := matching(ctx, "nothing"); e == nil || res != nil {
		t.Fatal("expected error", e, res)
	}
}

func matching(ctx Context, phrase string) (ret Result, err error) {
	match := &Noun{}
	words := strings.Fields(phrase)
	if bounds, e := ctx.GetBounds("", ""); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, bounds, Cursor{Words: words})
	}
	return
}

func matchingFilter(ctx Context, phrase, attr, class string) (ret Result, err error) {
	match := &Noun{Filters{&HasAttr{attr}, &HasClass{class}}}
	words := strings.Fields(phrase)
	if bounds, e := ctx.GetBounds("", ""); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, bounds, Cursor{Words: words})
	}
	return
}
