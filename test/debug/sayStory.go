package debug

import (
  "git.sr.ht/~ionous/iffy/affine"
  "git.sr.ht/~ionous/iffy/dl/core"
  "git.sr.ht/~ionous/iffy/dl/literal"
  "git.sr.ht/~ionous/iffy/jsn"
  "git.sr.ht/~ionous/iffy/rt"
  g "git.sr.ht/~ionous/iffy/rt/generic"
  "git.sr.ht/~ionous/iffy/rt/safe"
  "git.sr.ht/~ionous/iffy/test/testpat"
)

func SayIt(s string) rt.Execute {
  return &core.Say{T(s)}
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
  if a, e := safe.CheckVariable(run, numVar.String(), affine.Number); e != nil {
    err = e
  } else {
    n := a.Int()
    ret = g.BoolOf(n == op.Val)
  }
  return
}

func DetermineSay(i int) *core.CallPattern {
  return &core.CallPattern{
    Pattern: core.PatternName{Str: "say_me"},
    Arguments: core.NamedArgs(
      "num", &core.FromNum{I(i)}),
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
  Name:   "say_me",
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

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Bool: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Num: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Num: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Text: s} }

var SayHelloGoodbye = core.NewActivity(
  &core.ChooseAction{
    If: B(true),
    Do: core.MakeActivity(&core.Say{
      Text: T("hello"),
    }),
    Else: &core.ChooseNothingElse{
      core.MakeActivity(&core.Say{
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
                      "type": "say_text",
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
                            "type": "say_text",
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
