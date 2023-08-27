package list_test

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/rt/scope"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func newListTime(src []string, p testpat.Map) (ret rt.Runtime, vals *g.Record, err error) {
	var kinds testutil.Kinds
	type Locals struct{ Source []string }
	kinds.AddKinds((*Locals)(nil))
	locals := kinds.NewRecord("locals")
	lt := testpat.Runtime{
		Map: p,
		Runtime: testutil.Runtime{
			Kinds: &kinds,
			Chain: scope.MakeChain(
				scope.FromRecord(locals),
			),
		}}
	if e := locals.SetNamedField("source", g.StringsOf(src)); e != nil {
		err = e
	} else {
		ret = &lt
		vals = locals
	}
	return
}

// cmd to collect some text into a list of strings.
type Write struct {
	out  *[]string
	Text rt.TextEval
}

// Execute writes text to the runtime's current writer.
func (op *Write) Execute(run rt.Runtime) (err error) {
	if t, e := op.Text.GetText(run); e != nil {
		err = e
	} else {
		(*op.out) = append((*op.out), t.String())
	}
	return
}

func joinText(run rt.Runtime, op rt.TextListEval) (ret string) {
	if vs, e := op.GetTextList(run); e != nil {
		ret = e.Error()
	} else {
		ret = joinStrings(vs.Strings())
	}
	return
}

func joinStrings(vs []string) (ret string) {
	if len(vs) > 0 {
		ret = strings.Join(vs, ", ")
	} else {
		ret = "-"
	}
	return
}
