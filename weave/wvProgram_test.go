package weave_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"github.com/kr/pretty"
)

// grammar parsing doesn't do very many useful things modelling wise;
// so this just tests that it gets into the database.
func TestGrammar(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("b"),
		&eph.Directives{
			Name: `jump/skip/hop`,
			Directive: grammar.Directive{
				Lede: []string{"jump", "skip", "hop"},
				Scans: []grammar.ScannerMaker{
					&grammar.Action{Action: "jumping"},
				},
			},
		},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadGrammar(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		`b:jump/skip/hop:{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}`,
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
