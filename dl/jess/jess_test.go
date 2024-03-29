package jess_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/jesstest"
	"git.sr.ht/~ionous/tapestry/support/match"
)

func TestPhrases(t *testing.T) {
	var skipped int
	const at = -1
	for i, p := range jesstest.Phrases {
		if i != at && at >= 0 {
			continue
		}
		if str, ok := p.Test(); !ok {
			if len(str) > 0 {
				t.Log("skipped", str)
				skipped++
			}
		} else {
			// reset the dynamic noun pool every test
			known.nounPool = make(map[string]string)
			for _, n := range known.nouns {
				name := n.String()
				known.nounPool[name] = name
				known.nounPool["$"+name] = "things"
			}
			// request on logging
			q := jess.AddContext(&known, jess.LogMatches)
			// create the test helper
			m := jesstest.MakeMock(q, known.nounPool, known.verbs)
			// run the test:
			t.Logf("testing: %d %s", i, str)
			if !p.Verify(m.Generate(str)) {
				t.Logf("failed %d", i)
				t.Fail()
			}
		}
	}
	if skipped > 0 {
		t.Logf("skipped %d tests", skipped)
	}
}

type w struct {
	jesstest.Mock
	info
}

type info struct {
	kinds []string
	traits, fields,
	nouns, directions, verbNames match.SpanList
	nounPool map[string]string
	verbs    map[string]jesstest.MockVerb
}

func (n *info) GetContext() int {
	return 0
}
func (n *info) FindKind(ws match.Span, out *kindsOf.Kinds) (ret string, width int) {
	str := strings.ToLower(ws.String())
	for i, k := range n.kinds {
		if strings.HasPrefix(str, k) {
			if i&1 == 0 { // singular are the even numbers
				k = n.kinds[i+1]
			}
			ret, width = k, 1 // always in the tests
			if out != nil {
				// hacks for testing
				k := kindsOf.Kind
				switch ret {
				case "color":
					k = kindsOf.Aspect
				case "groups":
					k = kindsOf.Record
				case "storing":
					k = kindsOf.Action
				}
				*out = k
			}
			break
		}
	}
	return
}
func (n *info) FindTrait(ws match.Span) (string, int) {
	m, cnt := n.traits.FindPrefix(ws)
	return m.String(), cnt
}
func (n *info) FindField(ws match.Span) (string, int) {
	m, cnt := n.fields.FindPrefix(ws)
	return m.String(), cnt
}
func (n *info) FindNoun(ws match.Span, kind string) (ret string, width int) {
	var m match.Span
	switch kind {
	case jess.Directions:
		m, width = n.directions.FindPrefix(ws)
		ret = m.String()
	case jess.Verbs:
		m, width = n.verbNames.FindPrefix(ws)
		ret = m.String()
	default:
		str := ws.String()
		if noun, ok := n.nounPool[str]; ok {
			ret, width = noun, len(ws)
		} else {
			m, width = n.nouns.FindExactMatch(ws)
			ret = m.String()
		}
	}
	return
}

var known = info{
	kinds: []string{
		"kind", "kinds",
		"object", "objects",
		"door", "doors",
		"room", "rooms",
		"direction", "directions",
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
		"aspect", "aspects",
		"color", "color", // aspect; uses the singular
		"group", "groups", // record
		"storing", "storing", //
	},
	traits: match.PanicSpans(
		"closed",
		"open",
		"openable",
		"transparent",
		"fixed in place",
		"dark",
	),
	fields: match.PanicSpans(
		"description",
		"title",
		"age",
	),
	nouns: match.PanicSpans(
		"story",
		"message",
		"missive",
	),
	directions: match.PanicSpans(
		"north", "south", "east", "west",
	),
	verbs:     verbs,
	verbNames: panicVerbs(),
}

// reduce the keys into spans for matching
func panicVerbs() match.SpanList {
	vs := make([]string, len(verbs))
	for v := range verbs {
		vs = append(vs, v)
	}
	return match.PanicSpans(vs...)
}

// fix? maybe add "wearing" instead of carrying, to test implication better?
var verbs = jesstest.MockVerbs{
	"carrying": {
		Subject:  "actors",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  "not worn",
		Reversed: false, // (parent) is carrying (child)
	},
	"carried by": {
		Subject:  "actors",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  "not worn",
		Reversed: true, // (child) is carried by (parent)
	},
	"in": {
		Subject:   "containers",
		Alternate: "rooms", // alternate
		Object:    "things",
		Relation:  "whereabouts",
		Implies:   "not worn",
		Reversed:  true, // (child) is in (parent)
	},
	"on": {
		Subject:  "supporters",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  "not worn",
		Reversed: true, // (child) is on (parent)
	},
	"suspicious of": {
		Subject:  "actors",
		Object:   "actors",
		Relation: "suspicion",
		Reversed: false, // (parent) is suspicious of (child)
	},
}
