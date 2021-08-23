// primDetailsTemplate.js
// note: partial templates seem to COPY arguments,
// so they lose .. relative hierarchy.
'use strict';
module.exports = `
func {{Pascal name}}_Detailed_Optional_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte, err error) {
  var zero {{Pascal name}}
  if val.{{Pascal uses}} != zero.{{Pascal uses}} {
    ret, err = {{Pascal name}}_Detailed_Marshal(n, val)
  }
  return
}
func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, val *{{Pascal name}}) ([]byte, error) {
  return json.Marshal(jsonexp.{{Pascal uses}}{
{{#if (IsPositioned this)}}
    Id: val.At.Offset,
{{/if}}
    Type:  {{Pascal name}}_Type,
    Value: val.{{Pascal uses}},
  })
}

var {{Pascal name}}_Detailed_Optional_Unmarshal = {{Pascal name}}_Detailed_Unmarshal

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
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
