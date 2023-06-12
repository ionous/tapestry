package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// grammar parsing doesn't do very many useful things modelling wise;
// so this just tests that it gets into the database.
func TestGrammar(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("b"),
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
	} else if out, e := dt.readGrammar(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		`b:jump/skip/hop:{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}`,
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
