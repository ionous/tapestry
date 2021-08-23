// primCompactTemplate.js
'use strict';
// note: partial templates seem to COPY arguments,
// so they lose .. relative hierarchy.
// compact primitive values serialize to json primitives: string, bool, number
// there's a wrinkle that some primitives ( aka overrides ) are also stored in memory as golang primitives
// while others are stored in structs: ex. type Prim string vs. type Prim struct { Str string }
// and further that those latter can be stored in memory with string keys ex. "$SOME_CHOICE"
module.exports = `
func {{Pascal name}}_Compact_Optional_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte, err error) {
  var zero {{Pascal name}}
  if val.{{Pascal uses}} != zero.{{Pascal uses}} {
    ret, err = {{Pascal name}}_Compact_Marshal(n, val)
  }
  return
}
func {{Pascal name}}_Compact_Marshal(n jsonexp.Context, val *{{Pascal name}}) ([]byte, error) {
{{#if (Uses name 'str')}}
  var out string
  if str, ok:= composer.FindChoice(val, val.Str); !ok {
    out = val.Str
  } else {
    out = str
  }
  return json.Marshal(out)
{{else}}
  return json.Marshal(val.{{Pascal uses}})
{{/if}}
}


var {{Pascal name}}_Compact_Optional_Unmarshal = {{Pascal name}}_Compact_Unmarshal

func {{Pascal name}}_Compact_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.{{Pascal uses}}
  if len(b) > 0 {
    if e := json.Unmarshal(b, &msg); e != nil {
      err =  errutil.New({{Pascal name}}_Type, "-", e)
    }
  }
  if err == nil {
{{#if (IsPositioned this)~}}
    out.At = reader.Position{Source: n.Source(), Offset: msg.Id}
{{/if}}
    out.{{Pascal uses}} = msg.Value
  }
  return
}
`;
