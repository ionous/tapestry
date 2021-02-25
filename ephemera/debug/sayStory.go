package debug

import (
  "git.sr.ht/~ionous/iffy/affine"
  "git.sr.ht/~ionous/iffy/dl/core"
  "git.sr.ht/~ionous/iffy/dl/pattern"
  "git.sr.ht/~ionous/iffy/rt"
  g "git.sr.ht/~ionous/iffy/rt/generic"
  "git.sr.ht/~ionous/iffy/rt/safe"
  "git.sr.ht/~ionous/iffy/test/testpat"
)

func SayIt(s string) rt.Execute {
  return &core.Say{&core.Text{s}}
}

type MatchNumber struct {
  Val int
}

func (m MatchNumber) GetBool(run rt.Runtime) (ret g.Value, err error) {
  if a, e := safe.CheckVariable(run, "num", affine.Number); e != nil {
    err = e
  } else {
    n := a.Int()
    ret = g.BoolOf(n == m.Val)
  }
  return
}

func DetermineSay(i int) *pattern.Determine {
  return &pattern.Determine{
    Pattern: "say_me",
    Arguments: core.NamedArgs(
      "num", &core.FromNum{
        &core.Number{float64(i)},
      }),
  }
}

type SayMe struct {
  Num float64
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
    If: &core.Bool{true},
    Do: core.MakeActivity(&core.Say{
      Text: &core.Text{"hello"},
    }),
    Else: &core.ChooseNothingElse{
      core.MakeActivity(&core.Say{
        Text: &core.Text{"goodbye"},
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
