// runTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} {{desc.short}}
type {{Pascal name}} struct {
{{#each (ParamsOf this)}}
  {{Pascal @key}} {{TypeOf this}} \`if:"
  {{~#if (IsInternal label)}}internal{{else}}label={{LabelOf @key this @index}}{{/if}}
  {{~#if optional}},optional{{/if}}\
  {{~#if (OverrideOf type)}},type={{type}}{{/if}}"\`
{{/each}}
}

{{>spec spec=this~}}
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
const {{Pascal ../name}}_{{Pascal @key}}= "{{@key}}";
{{~/unless}}{{~/each}}
{{#if ../marshal}}
{{>sig sig=this}}

func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte, err error) {
  var fields jsonexp.Fields
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
  if b, e:= {{ScopeOf type}}{{Pascal type}}_Detailed{{#if (OverrideOf type)}}_Override{{/if}}{{ModOf this}}_Marshal(n, &val.{{Pascal @key}}); e!= nil {
    err = errutil.Append(err, e)
  } else if len(b) > 0 {
      fields[{{Pascal ../name}}_{{Pascal @key}}]= b
  }{{/unless}}{{/each}}
  if err== nil {
    ret, err= json.Marshal(jsonexp.Flow{
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
  if e:= json.Unmarshal(b, &msg); e!= nil {
    err = e
  } {{#each (ParamsOf this)}}{{#unless (IsInternal label)~}}
  else if e:= {{ScopeOf type}}{{Pascal type}}_Detailed{{#if (OverrideOf type)}}_Override{{/if}}{{ModOf this}}_Unmarshal(n, msg.Fields[{{Pascal ../name}}_{{Pascal @key}}], &out.{{Pascal @key}}); e!= nil {
    err = e
  } {{/unless}}{{/each~}}
{{~#if (IsPositioned this)~}}
  else {
      out.At = reader.Position{Source:n.Source, Offset: msg.Id}
  }{{/if}}
  return
}
{{/if}}{{/with}}
{{#if marshal}}
{{>repeat repeat=this}}
{{>optional optional=this}}
{{/if}}
`;
