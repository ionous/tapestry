package debug

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

func SayIt(s string) []rt.Execute {
	return []rt.Execute{&format.PrintText{Text: T(s)}}
}

type MatchNumber struct {
	Val int
}

// num_value, a type of flow.
// var Zt_MatchNumber = typeinfo.Flow{
// 	Name: "match_number",
// 	Lede: "match",
// 	Terms: []typeinfo.Term{{
// 		Name: "val",
// 		Type: &prim.Zt_Number,
// 	}},
// 	Slots: []*typeinfo.Slot{
// 		&rtti.Zt_BoolEval,
// 	},
// }

func (op *MatchNumber) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := run.GetField(meta.Variables, "num"); e != nil {
		err = nil
	} else if safe.Check(v, affine.Num); e != nil {
		err = e
	} else {
		ret = rt.BoolOf(v.Int() == op.Val)
	}
	return
}

func DetermineSay(i int) *call.CallPattern {
	return &call.CallPattern{
		PatternName: "say me",
		Arguments: []call.Arg{{
			Name:  "num",
			Value: &call.FromNum{Value: I(i)},
		}},
	}
}

type SayMe struct {
	Num float64
}

// func (op *SayMe) Marshal(m jsn.Marshaler) (err error) {
// 	if err = m.MarshalBlock(MakeFlow(op)); err == nil {
// 		e0 := m.MarshalKey("", "")
// 		if e0 == nil {
// 			e0 = m.MarshalValue("", &op.Num)
// 		}
// 		if e0 != nil && e0 != jsn.Missing {
// 			m.Error(e0)
// 		}
// 		m.EndBlock()
// 	}
// 	return
// }

// the rules defined last run first
var SayPattern = testpat.Pattern{
	Name:   "say me",
	Labels: []string{"num"},
	Fields: []rt.Field{
		{Name: "num", Affinity: "number", Type: ""},
	},
	Rules: []rt.Rule{{
		Name: "default", Exe: SayIt("Not between 1 and 3."),
	}, {
		Name: "3b", Exe: []rt.Execute{&logic.ChooseBranch{
			Condition: &MatchNumber{3}, Exe: SayIt("San!")},
		},
	}, {
		Name: "3a", Exe: []rt.Execute{&logic.ChooseBranch{
			Condition: &MatchNumber{3}, Exe: SayIt("Three!")},
		},
	}, {
		Name: "2", Exe: []rt.Execute{&logic.ChooseBranch{
			Condition: &MatchNumber{2}, Exe: SayIt("Two!")},
		},
	}, {
		Name: "1", Exe: []rt.Execute{&logic.ChooseBranch{
			Condition: &MatchNumber{1}, Exe: SayIt("One!")},
		},
	}},
}
