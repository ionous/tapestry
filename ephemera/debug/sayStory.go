package debug

import (
  "git.sr.ht/~ionous/iffy/affine"
  "git.sr.ht/~ionous/iffy/dl/core"
  "git.sr.ht/~ionous/iffy/dl/value"
  "git.sr.ht/~ionous/iffy/jsn"
  "git.sr.ht/~ionous/iffy/rt"
  g "git.sr.ht/~ionous/iffy/rt/generic"
  "git.sr.ht/~ionous/iffy/rt/safe"
  "git.sr.ht/~ionous/iffy/test/testpat"
)

func SayIt(s string) rt.Execute {
  return &core.Say{&core.TextValue{s}}
}

type MatchNumber struct {
  Val int
}

func (op *MatchNumber) Marshal(m jsn.Marshaler) (err error) {
  if err = m.MarshalBlock(jsn.MakeFlow("match", "", op)); err == nil {
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
  if a, e := safe.CheckVariable(run, numVar, affine.Number); e != nil {
    err = e
  } else {
    n := a.Int()
    ret = g.BoolOf(n == op.Val)
  }
  return
}

func DetermineSay(i int) *core.CallPattern {
  return &core.CallPattern{
    Pattern: value.PatternName{Str: "say_me"},
    Arguments: core.NamedArgs(
      "num", &core.FromNum{
        &core.NumValue{float64(i)},
      }),
  }
}

type SayMe struct {
  Num float64
}

func (op *SayMe) Marshal(m jsn.Marshaler) (err error) {
  if err = m.MarshalBlock(jsn.MakeFlow("say_me", "", op)); err == nil {
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

var SayHelloGoodbye = core.NewActivity(
  &core.ChooseAction{
    If: &core.BoolValue{true},
    Do: core.MakeActivity(&core.Say{
      Text: &core.TextValue{"hello"},
    }),
    Else: &core.ChooseNothingElse{
      core.MakeActivity(&core.Say{
        Text: &core.TextValue{"goodbye"},
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
