package parser_test

import (
	"strings"
	"testing"

	. "git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/parser/ident"
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
	"github.com/ionous/sliceOf"
)

// MyObject - along with MyNoun this provides an example of mapping some application defined object to a parser.NounInstance.
type MyObject struct {
	Id         ident.Id
	Names      []string
	Classes    []string
	Attributes []string
}

func (m *MyObject) String() string {
	return m.Id.String()
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

func (m MyBounds) GetPlayerBounds(n string) (ret Bounds, err error) {
	switch n {
	case "":
		ret = m.SearchBounds
	case "self":
		ret = m.SelfBounds
	default:
		err = errutil.New("unknown bounds", n)
	}
	return
}

func (m MyBounds) GetObjectBounds(ident.Id) (Bounds, error) {
	return m.SearchBounds, nil
}

func (m MyBounds) IsPlural(word string) bool {
	return word != inflect.Singularize(word)
}

func (m MyBounds) SelfBounds(v NounVisitor) (ret bool) {
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

func (adapt MyNoun) Id() ident.Id {
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
		&MyObject{Id: ident.IdOf("a"), Names: sliceOf.String("unique")},
		//
		&MyObject{Id: ident.IdOf("b"), Names: strings.Fields("exact")},
		&MyObject{Id: ident.IdOf("c"), Names: strings.Fields("exact match")},
		//
		&MyObject{Id: ident.IdOf("d"), Names: strings.Fields("inexact match")},
		&MyObject{Id: ident.IdOf("e"), Names: strings.Fields("inexact conflict")},
		//
		&MyObject{Id: ident.IdOf("f"),
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
			Classes:    strings.Fields("class"),
		},
		&MyObject{Id: ident.IdOf("g"),
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
		},
		&MyObject{Id: ident.IdOf("h"),
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
	if bounds, e := ctx.GetPlayerBounds(""); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, bounds, Cursor{Words: words})
	}
	return
}

func matchingFilter(ctx Context, phrase, attr, class string) (ret Result, err error) {
	match := &Noun{Filters{&HasAttr{attr}, &HasClass{class}}}
	words := strings.Fields(phrase)
	if bounds, e := ctx.GetPlayerBounds(""); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, bounds, Cursor{Words: words})
	}
	return
}
