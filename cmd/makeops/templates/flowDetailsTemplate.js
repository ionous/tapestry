// flowDetails.js
'use strict';
module.exports = `
func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte,err error) {
{{#unless (ParamsOf this)}}
  var fields jsonexp.Fields
{{else}}
  fields := make(jsonexp.Fields)
{{~/unless}}
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
  if b,e := {{ScopeOf type}}{{Pascal type}}_Detailed{{#if (OverrideOf type)}}_Override{{/if}}{{ModOf this}}_Marshal(n, &val.{{Pascal @key}}); e != nil {
    err = errutil.Append(err, e)
  } else if len(b) > 0 {
      fields[{{Pascal ../name}}_{{Pascal @key}}]= b
  }{{/unless}}{{/each}}
  if err == nil {
    ret, err = json.Marshal(jsonexp.Flow{
{{~#if (IsPositioned this)}}
    Id: val.At.Offset,{{/if}}
      Type:  Type_{{Pascal name}},
      Fields: fields,
    })
  }
  return
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.Flow
  if e := json.Unmarshal(b, &msg); e != nil {
    err = errutil.New(Type_{{Pascal name}}, "-", e)
  } {{#each (ParamsOf this)}}{{#unless (IsInternal label)~}}
  else if e := {{ScopeOf type}}{{Pascal type}}_Detailed{{#if (OverrideOf type)}}_Override{{/if}}{{ModOf this}}_Unmarshal(n, msg.Fields[{{Pascal ../name}}_{{Pascal @key}}], &out.{{Pascal @key}}); e != nil {
    err = errutil.New(Type_{{Pascal ../name}} + "." + {{Pascal ../name}}_{{Pascal @key}}, "-", e)
  } {{/unless}}{{/each~}}
{{~#if (IsPositioned this)~}}
  else {
      out.At = reader.Position{Source: n.Source(), Offset: msg.Id}
  }{{/if}}
  return
}
`
