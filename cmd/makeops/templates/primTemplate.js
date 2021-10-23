// strTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  {{Pascal uses}} {{#if (Uses name 'num')}}float64{{else}}string{{/if}}
}

{{#if (Uses name 'str')}}
func (op *{{Pascal name}}) String() string {
  return op.Str
}
{{/if}}

{{#each (Choices @this)}}
const {{Pascal ../name}}_{{Pascal this.token}}= "{{this.token}}";
{{/each}}

{{>spec}}
{{#if ../marshal}}
{{>sig}}
{{#if (Unboxed name)}}
{{>repeat name=(Pascal name) mod="_Unboxed" el=(Unboxed name)}}

func {{Pascal name}}_Unboxed_Optional_Marshal(m jsn.Marshaler, val *{{Unboxed name}}) (err error) {
  var zero {{Unboxed name}}
  if enc := m.IsEncoding(); !enc || *val != zero {
    err = {{Pascal name}}_Unboxed_Marshal(m, val)
  }
  return
}

func {{Pascal name}}_Unboxed_Marshal(m jsn.Marshaler, val *{{Unboxed name}}) error {
  return m.MarshalValue({{Pascal name}}_Type, jsn.Box{{Pascal (Unboxed name)}}(val))
}
{{/if}}

func {{Pascal name}}_Optional_Marshal(m jsn.Marshaler, val *{{Pascal name}}) (err error) {
  var zero {{Pascal name}}
  if enc := m.IsEncoding(); !enc || val.{{Pascal uses}} != zero.{{Pascal uses}} {
     err = {{Pascal name}}_Marshal(m, val)
  }
  return
}

func {{Pascal name}}_Marshal(m jsn.Marshaler, val *{{Pascal name}}) (err error) {
{{#if (IsCustom name)}}
  if fn, ok := m.CustomizedMarshal({{Pascal name}}_Type); ok {
    err = fn(m, val)
  } else {
    err = {{Pascal name}}_DefaultMarshal(m, val)
  }
}
func {{Pascal name}}_DefaultMarshal(m jsn.Marshaler, val *{{Pascal name}}) error {
{{/if}}
{{#if (IsPositioned this)}}
  m.SetCursor(val.At.Offset)
{{/if}}
  return m.MarshalValue({{Pascal name}}_Type, {{#if (IsEnumerated this)~}}
    jsn.MakeEnum(val, &val.Str){{else~}}
    &val.{{Pascal uses}}{{/if~}}
  )
}

{{/if}}
{{>repeat name=(Pascal name) el=(Pascal name)}}
{{/with}}
`;
