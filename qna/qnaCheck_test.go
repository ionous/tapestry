package qna_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestCheck(t *testing.T) {
	prog := &qna.CheckOutput{
		Name:   t.Name(),
		Expect: "hello",
		Test: []rt.Execute{
			&core.ChooseAction{
				If: &literal.BoolValue{Bool: true},
				Does: core.MakeActivity(&core.Say{
					Text: &literal.TextValue{Text: "hello"},
				}),
				Else: &core.ChooseNothingElse{
					Does: core.MakeActivity(&core.Say{
						Text: &literal.TextValue{Text: "goodbye"},
					})},
			}},
	}
	if e := runTest(prog); e != nil {
		t.Fatal(e)
	}
}

func runTest(prog *qna.CheckOutput) (err error) {
	var run checkTester
	run.SetWriter(print.NewAutoWriter(writer.NewStdout()))
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
