package qna_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/qna"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

func TestCheck(t *testing.T) {
	prog := &qna.CheckOutput{
		Name:   t.Name(),
		Expect: "hello",
		Test: &core.Activity{Exe: []rt.Execute{
			&core.ChooseAction{
				If: &literal.BoolValue{true},
				Do: core.MakeActivity(&core.Say{
					Text: &literal.TextValue{"hello"},
				}),
				Else: &core.ChooseNothingElse{
					Do: core.MakeActivity(&core.Say{
						Text: &literal.TextValue{"goodbye"},
					})},
			}}},
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
