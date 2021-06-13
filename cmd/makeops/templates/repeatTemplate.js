// repeatsPartial.js
'use strict';
module.exports =`{{#if repeats}}{{#with type}}
func {{Pascal name}}_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]{{Pascal name}}) (ret []byte,err error) {
  var msgs []json.RawMessage
  msgs= make([]json.RawMessage, len(*vals))
  for i, el:= range *vals {
    if b,e := {{Pascal name}}_Detailed_Marshal(n, &el); e!= nil {
      err =  errutil.New(Type_{{Pascal name}}, "at", i, "-", e)
      break
    } else {
      msgs[i]= b
    }
  }
  if err== nil {
    ret, err= json.Marshal(msgs)
  }
  return
}

func {{Pascal name}}_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]{{Pascal name}}) (err error) {
  var vals []{{Pascal name}}
  if len(b) > 0 { // generated code collapses optional and empty.
    var msgs []json.RawMessage
    if e := json.Unmarshal(b, &msgs); e!= nil  {
      err =  errutil.New(Type_{{Pascal name}}, "-", e)
    } else {
      vals = make([]{{Pascal name}}, len(msgs))
      for i, msg:= range msgs {
        if e := {{Pascal name}}_Detailed_Unmarshal(n, msg, &vals[i]); e!= nil {
          err =  errutil.New(Type_{{Pascal name}}, "at", i, "-", e)
          break
        }
      }
    }
  }
  if err== nil {
    *out= vals
  }
  return
}
{{/with}}{{/if}}`;
