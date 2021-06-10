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
    ret, err= json.Marshal(jsonexp.Flow{
{{#if (IsPositioned this)}}
      Id: op.At.Offset,
{{/if}}      Type:  Type_{{Pascal name}},
      Value: map[string]json.RawMessage{
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
        "{{@key}}": json{{Pascal @key}},
{{~/unless}}{{~/each}}
      },
    })
  }
  return
}

{{#each (ParamsOf this)}}{{#unless (IsInternal label)}}
func (op *{{Pascal ../name}}) MarshalJSON{{Pascal @key}}() (ret []byte, err error) {
{{#if repeats}}
  ret, err= json.Marshal( op.{{Pascal @key}} )
{{else if (IsBool this)}}
  // bool override
  var str string
  if op.{{Pascal @key}} {
    str= value.Bool_True
  } else {
     str= value.Bool_False
  }
  m:= value.Bool{ str }
  ret, err= m.MarshalJSON()
{{else if (Override this)}}
  // type override
   m:= {{ScopeOf type}}{{Pascal type}} { op.{{Pascal @key}} }
   ret, err= m.MarshalJSON()
{{else if (Uses type 'slot')}}
  if v, e:= op.{{Pascal @key}}.(json.Marshaler).MarshalJSON(); e!= nil {
    err= e
  } else {
    ret, err= json.Marshal(jsonexp.Slot{
      Type:  {{ScopeOf type}}Type_{{Pascal type}},
      Value: v,
    })
  }
{{else}}
  ret, err= op.{{Pascal @key}}.MarshalJSON()
{{/if}}
  return
}
{{/unless}}{{/each}}{{/if}}
{{/with}}
`;
