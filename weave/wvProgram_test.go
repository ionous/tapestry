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
			Directive: grammar.Directive{
				Name: `jump/skip/hop`,
				Scans: []grammar.ScannerMaker{
					&grammar.Words{Words: []string{"jump", "skip", "hop"}},
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
		`b:{"Interpret name:with:":["jump/skip/hop",[{"One word:":["jump","skip","hop"]},{"Action:":"jumping"}]]}`,
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(out))
		t.Fatal(diff)
	}
}
