package debug

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

func SayIt(s string) []rt.Execute {
	return []rt.Execute{&core.PrintText{Text: T(s)}}
}

type MatchNumber struct {
	Val int
}

func (op *MatchNumber) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(MakeFlow(op)); err == nil {
		e0 := m.MarshalKey("", "")
		if e0 == nil {
			e0 = m.MarshalValue("", &op.Val)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(e0)
		}
		m.EndBlock()
	}
	return
}

func (op *MatchNumber) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := run.GetField(meta.Variables, "num"); e != nil {
		err = nil
	} else if safe.Check(v, affine.Number); e != nil {
		err = e
	} else {
		ret = g.BoolOf(v.Int() == op.Val)
	}
	return
}

func DetermineSay(i int) *assign.CallPattern {
	return &assign.CallPattern{
		PatternName: "say me",
		Arguments: []assign.Arg{{
			Name:  "num",
			Value: &assign.FromNumber{Value: I(i)},
		}},
	}
}

type SayMe struct {
	Num float64
}

func (op *SayMe) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(MakeFlow(op)); err == nil {
		e0 := m.MarshalKey("", "")
		if e0 == nil {
			e0 = m.MarshalValue("", &op.Num)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(e0)
		}
		m.EndBlock()
	}
	return
}

var SayPattern = testpat.Pattern{
	Name:   "say me",
	Labels: []string{"num"},
	Fields: []g.Field{
		{Name: "num", Affinity: "number", Type: ""},
	},
	Rules: []rt.Rule{
		{Name: "default", Execute: SayIt("Not between 1 and 3.")},
		{Name: "3b", Filter: &MatchNumber{3}, Execute: SayIt("San!")},
		{Name: "3a", Filter: &MatchNumber{3}, Execute: SayIt("Three!")},
		{Name: "2", Filter: &MatchNumber{2}, Execute: SayIt("Two!")},
		{Name: "1", Filter: &MatchNumber{1}, Execute: SayIt("One!")},
	},
}

var SayHelloGoodbye = core.MakeActivity(
	&core.ChooseBranch{
		If: B(true),
		Exe: core.MakeActivity(&core.PrintText{
			Text: T("hello"),
		}),
		Else: &core.ChooseNothingElse{Exe: core.MakeActivity(&core.PrintText{
			Text: T("goodbye"),
		}),
		},
	})

var SayHelloGoodbyeData = `{
  "type": "activity",
  "value": {
    "$EXE": [{
        "type": "execute",
        "value": {
          "type": "choose_action",
          "value": {
            "$DO": {
              "type": "activity",
              "value": {
                "$EXE": [{
                    "type": "execute",
                    "value": {
                      "type": "say",
                      "value": {
                        "$TEXT": {
                          "type": "text_eval",
                          "value": {
                            "type": "text_value",
                            "value": {
                              "$TEXT": {
                                "type": "text",
                                "value": "hello"
                              }}}}}}}]}},
            "$ELSE": {
              "type": "brancher",
              "value": {
                "type": "choose_nothing_else",
                "value": {
                  "$DO": {
                    "type": "activity",
                    "value": {
                      "$EXE": [
                        {
                          "type": "execute",
                          "value": {
                            "type": "say",
                            "value": {
                              "$TEXT": {
                                "type": "text_eval",
                                "value": {
                                  "type": "text_value",
                                  "value": {
                                    "$TEXT": {
                                      "type": "text",
                                      "value": "goodbye"
                                    }}}}}}}]}}}}},
            "$IF": {
              "type": "bool_eval",
              "value": {
                "type": "bool_value",
                "value": {
                  "$BOOL": {
                    "type": "bool",
                    "value": "$TRUE"
                  }}}}}}}]}}
`
