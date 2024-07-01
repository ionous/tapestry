package cmdcheck_test

import (
	"os"
	"testing"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cmdcheck"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"git.sr.ht/~ionous/tapestry/web/markup"
)

func TestCheck(t *testing.T) {
	prog := query.CheckData{
		Name:   t.Name(),
		Expect: "hello",
		Test: []rt.Execute{
			&logic.ChooseBranch{
				Condition: &literal.BoolValue{Value: true},
				Exe: []rt.Execute{&format.PrintText{
					Text: &literal.TextValue{Value: "hello"},
				}},
				Else: &logic.ChooseNothingElse{
					Exe: []rt.Execute{&format.PrintText{
						Text: &literal.TextValue{Value: "goodbye"},
					}}},
			}},
	}
	if e := runTest(prog); e != nil {
		t.Fatal(e)
	}
}

func runTest(prog query.CheckData) (err error) {
	var run checkTester
	run.SetWriter(print.NewLineSentences(markup.ToText(os.Stdout)))
	return cmdcheck.RunTest(&run, prog)
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type checkTester struct {
	baseRuntime
	writer.Sink
}

func (c *checkTester) ActivateDomain(string) (_ error) { return }
