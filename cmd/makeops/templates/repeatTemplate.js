// repeatsPartial.js
// note: the non-repeating callbacks used by these functions are expected to be defined elsewhere
// ex. in flow, or prim, etc.
'use strict';
module.exports =`
func {{name}}{{mod}}_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]{{el}}) ([]byte, error) {
  return {{name}}{{mod}}_Repeats_Marshal(n, vals, {{name}}{{mod}}_Compact_Marshal)
}
func {{name}}{{mod}}_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]{{el}}) ([]byte, error) {
  return {{name}}{{mod}}_Repeats_Marshal(n, vals, {{name}}{{mod}}_Detailed_Marshal)
}
func {{name}}{{mod}}_Repeats_Marshal(n jsonexp.Context, vals *[]{{el}}, marshEl func(jsonexp.Context, *{{el}}) ([]byte, error)) (ret []byte, err error) {
  var msgs []json.RawMessage
  if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
    msgs = make([]json.RawMessage, cnt)
    for i, el := range *vals {
      if b,e := marshEl(n, &el); e != nil {
        err =  errutil.New({{name}}_Type, "at", i, "-", e)
        break
      } else {
        msgs[i]= b
      }
    }
  }
  if err == nil {
    ret, err = json.Marshal(msgs)
  }
  return
}

func {{name}}{{mod}}_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]{{el}}) error {
  return {{name}}{{mod}}_Repeats_Unmarshal(n, b, out, {{name}}{{mod}}_Compact_Unmarshal)
}
func {{name}}{{mod}}_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]{{el}}) error {
  return {{name}}{{mod}}_Repeats_Unmarshal(n, b, out, {{name}}{{mod}}_Detailed_Unmarshal)
}
func {{name}}{{mod}}_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]{{el}}, unmarshEl func(jsonexp.Context, []byte, *{{el}}) error) (err error) {
  var vals []{{el}}
  if len(b) > 0 { // generated code collapses optional and empty.
    var msgs []json.RawMessage
    if e := json.Unmarshal(b, &msgs); e != nil  {
      err =  errutil.New({{name}}_Type, "-", e)
    } else {
      vals = make([]{{el}}, len(msgs))
      for i, msg := range msgs {
        if e := unmarshEl(n, msg, &vals[i]); e != nil {
          err =  errutil.New({{name}}_Type, "at", i, "-", e)
          break
        }
      }
    }
  }
  if err == nil {
    *out= vals
  }
  return
}
`;
