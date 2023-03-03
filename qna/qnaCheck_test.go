package qna_test

import (
	"os"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"git.sr.ht/~ionous/tapestry/web/markup"
)

func TestCheck(t *testing.T) {
	prog := &qna.CheckOutput{
		Name:   t.Name(),
		Expect: "hello",
		Test: []rt.Execute{
			&core.ChooseAction{
				If: &literal.BoolValue{Value: true},
				Does: core.MakeActivity(&core.PrintText{
					Text: &literal.TextValue{Value: "hello"},
				}),
				Else: &core.ChooseNothingElse{
					Does: core.MakeActivity(&core.PrintText{
						Text: &literal.TextValue{Value: "goodbye"},
					})},
			}},
	}
	if e := runTest(prog); e != nil {
		t.Fatal(e)
	}
}

func runTest(prog *qna.CheckOutput) (err error) {
	var run checkTester
	run.SetWriter(print.NewLineSentences(markup.ToText(os.Stdout)))
	return prog.RunTest(&run)
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type checkTester struct {
	baseRuntime
	writer.Sink
}

func (c *checkTester) ActivateDomain(string) (prev string, err error) { return }
