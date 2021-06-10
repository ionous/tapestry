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
  {{~#if (Override this)}},type={{Override this}}{{/if}}"\`
{{/each}}
}

{{>spec spec=this}}

{{#if ../marshal}}
func (op *{{Pascal name}}) MarshalJSON() (ret []byte, err error) {
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
  if json{{Pascal @key}}, e:= op.MarshalJSON{{Pascal @key}}(); e!= nil {
    err = e
  } else {{/unless}}{{/each~}}
  {
    ret, err= json.Marshal(map[string]interface{}{
{{#if (IsPositioned this)}}
      "id": op.At.Offset,
{{/if}}      "type":  "{{name}}",
      "value": map[string]json.RawMessage{
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
        "{{@key}}": json{{Pascal @key}},
{{~/unless}}{{~/each}}
      },
    })
  }
  return
}

{{#each (ParamsOf this)}}{{#unless (IsInternal label)}}
func (op *{{Pascal ../name}}) MarshalJSON{{Pascal @key}}() ([]byte, error) {
{{#if repeats}}
  return json.Marshal( op.{{Pascal @key}} )
{{else if (IsBool this)}}
  // bool override
  var str string
  if op.{{Pascal @key}} {
    str= value.Bool_True
  } else {
     str= value.Bool_False
  }
  m:= value.Bool{ str }
  return m.MarshalJSON()
{{else if (Override this)}}
  // type override
   m:= {{OriginalTypeOf this}} { op.{{Pascal @key}} }
   return m.MarshalJSON()
{{else if (Uses type 'slot')}}
  m:= op.{{Pascal @key}}.(json.Marshaler)
  return m.MarshalJSON()
{{else}}
  return op.{{Pascal @key}}.MarshalJSON()
{{/if}}
}
{{/unless}}{{/each}}{{/if}}
{{/with}}
`;
