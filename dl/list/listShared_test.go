package list_test

import (
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testpat"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

func newListTime(src []string, p pattern.Map) (ret rt.Runtime, vals *g.Record, err error) {
	var kinds testutil.Kinds
	type Values struct{ Source []string }
	kinds.AddKinds((*Values)(nil))
	values := kinds.New("values")
	lt := testpat.Runtime{
		p,
		testutil.Runtime{
			Kinds: &kinds,
			Stack: []rt.Scope{
				g.RecordOf(values),
			},
		}}
	if e := values.SetNamedField("source", g.StringsOf(src)); e != nil {
		err = e
	} else {
		ret = &lt
		vals = values
	}
	return
}

func B(i bool) rt.BoolEval     { return &core.Bool{i} }
func I(i int) rt.NumberEval    { return &core.Number{float64(i)} }
func T(i string) rt.TextEval   { return &core.Text{i} }
func V(i string) *core.Var     { return &core.Var{Name: i} }
func N(n string) core.Variable { return core.Variable{Str: n} }

func FromTs(vs []string) (ret rt.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{&core.Text{vs[0]}}
	} else {
		ret = &core.FromTexts{&core.Texts{vs}}
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

func getNum(run rt.Runtime, op rt.NumberEval) (ret string) {
	if v, e := op.GetNumber(run); e != nil {
		ret = e.Error()
	} else {
		ret = strconv.FormatFloat(v.Float(), 'g', -1, 64)
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
