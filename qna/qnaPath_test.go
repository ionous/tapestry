package qna_test

import (
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"testing"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

func TestEventPath(t *testing.T) {
	type EventValues struct {
		Objects []string
	}
	type Things struct {
	}
	type Devices struct {
		Things
	}

	var kinds testutil.Kinds
	var objects testutil.Objects
	kinds.AddKinds((*Things)(nil), (*Devices)(nil), (*EventValues)(nil))
	objects.AddObjects(kinds.Kind("things"), "table", "apple")
	objects.AddObjects(kinds.Kind("devices"), "pen")
	// create a "values" containing the field "objects" with a list of all object names
	values := kinds.NewRecord("event_values", "objects", objects.Names())
	//
	run := &eventPathRt{
		testpat.Runtime{
			Runtime: testutil.Runtime{
				Kinds:     &kinds,
				ObjectMap: objects,
				Stack: []rt.Scope{
					scope.FromRecord(values),
				},
			},
		},
	}

	// pen is a child of table
	if path, e := qna.BuildPath(run, "event", []string{"pen", "table"}, nil); e != nil {
		t.Fatal(e)
	} else {
		var got []string
		for _, p := range path {
			for _, rule := range p {
				got = append(got, rule.Name)
			}
		}
		if diff := pretty.Diff(got, []string{
			"pen-1-Prefix", "pen-0-Prefix", "devices-1-Prefix", "devices-0-Prefix", "things-1-Prefix", "things-0-Prefix",
			"table-1-Prefix", "table-0-Prefix", "things-1-Prefix", "things-0-Prefix",
			"pen-1-Infix", "pen-0-Infix", "devices-1-Infix", "devices-0-Infix", "things-1-Infix", "things-0-Infix",
			// this is the most questionable row.
			// do we really want to visit rules for the table if the pen is the object in question?
			"table-1-Infix", "table-0-Infix", "things-1-Infix", "things-0-Infix",
			"pen-1-Postfix", "pen-0-Postfix", "devices-1-Postfix", "devices-0-Postfix", "things-1-Postfix", "things-0-Postfix",
			"table-1-Postfix", "table-0-Postfix", "things-1-Postfix", "things-0-Postfix",
			"pen-1-After", "pen-0-After", "devices-1-After", "devices-0-After", "things-1-After", "things-0-After",
			"table-1-After", "table-0-After", "things-1-After", "things-0-After"}); len(diff) > 0 {
			t.Fatal(got)
		}
	}
}

type eventPathRt struct {
	testpat.Runtime
}

// we can assume rues gives us rules ordered by flags, with rules declared last listed first
// ( order by phase, mr.rowid )
func (ep *eventPathRt) GetRules(pattern, target string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	// generate four sets of rules:
	for p := 0; p < 4; p++ {
		flags := rt.Flags(1 << p)
		if pflags != nil {
			*pflags |= flags
		}
		// each with two layers to simulate "recently" declared and "oldest" declared
		ret = append(ret,
			rt.Rule{Name: target + "-1-" + flags.String(), RawFlags: float64(flags)},
			rt.Rule{Name: target + "-0-" + flags.String(), RawFlags: float64(flags)})
	}
	return
}
