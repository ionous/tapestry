'use strict';
module.exports = `
func {{Pascal name}}_Override_{{fmt}}_Optional_Marshal(n jsonexp.Context, val *{{OverrideOf name}}) (ret []byte, err error) {
  var zero {{OverrideOf name}}
  if *val != zero {
    ret, err = {{Pascal name}}_Override_{{fmt}}_Marshal(n, val)
  }
  return
}
func {{Pascal name}}_Override_{{fmt}}_Marshal(n jsonexp.Context, val *{{OverrideOf name}}) ([]byte, error) {
{{#unless (IsBool name)}}
  return {{Pascal name}}_{{fmt}}_Marshal(n, &{{Pascal name}}{*val})
{{else}}
  var out string
  if *val {
    out = Bool_True
  } else {
    out = Bool_False
  }
  return {{Pascal name}}_{{fmt}}_Marshal(n, &{{Pascal name}}{out})
{{/unless}}
}

var {{Pascal name}}_Override_{{fmt}}_Optional_Unmarshal = {{Pascal name}}_Override_{{fmt}}_Unmarshal

func {{Pascal name}}_Override_{{fmt}}_Unmarshal(n jsonexp.Context, b []byte, out *{{OverrideOf name}}) (err error) {
  if len(b) > 0 {
    var msg {{Pascal name}}
    if e := {{Pascal name}}_{{fmt}}_Unmarshal(n, b, &msg); e != nil {
      err = errutil.New({{Pascal name}}_Type, "-", e)
    } else {
  {{#unless (IsBool name)}}
      *out= msg.{{Pascal uses}}
  {{else}}
      *out = msg.{{Pascal uses}} == Bool_True
  {{/unless}}
    }
  }
  return
}
`;
