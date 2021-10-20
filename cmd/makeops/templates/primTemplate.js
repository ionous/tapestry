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

func {{Pascal name}}_Unboxed_Optional_Marshal(n jsn.Marshaler, val *{{Unboxed name}}) {
  var zero {{Unboxed name}}
  if enc := n.IsEncoding(); !enc || *val != zero {
    {{Pascal name}}_Unboxed_Marshal(n, val)
  }
}

func {{Pascal name}}_Unboxed_Marshal(n jsn.Marshaler, val *{{Unboxed name}}) {
  n.MarshalValue({{Pascal name}}_Type, jsn.Box{{Pascal (Unboxed name)}}(val))
}
{{/if}}

func {{Pascal name}}_Optional_Marshal(n jsn.Marshaler, val *{{Pascal name}}) {
  var zero {{Pascal name}}
  if enc := n.IsEncoding(); !enc || val.{{Pascal uses}} != zero.{{Pascal uses}} {
     {{Pascal name}}_Marshal(n, val)
  }
}

func {{Pascal name}}_Marshal(n jsn.Marshaler, val *{{Pascal name}}) {
{{#if (IsCustom name)}}
  if fn, ok := n.CustomizedMarshal({{Pascal name}}_Type); ok {
    fn(n, val)
  } else {
    {{Pascal name}}_DefaultMarshal(n, val)
  }
}
func {{Pascal name}}_DefaultMarshal(n jsn.Marshaler, val *{{Pascal name}}) {
{{/if}}
{{#if (IsPositioned this)}}
  n.SetCursor(val.At.Offset)
{{/if}}
  n.MarshalValue({{Pascal name}}_Type, {{#if (IsEnumerated this)~}}
    jsn.MakeEnum(val, &val.Str){{else~}}
    &val.{{Pascal uses}}{{/if~}}
  )
}

{{/if}}
{{>repeat name=(Pascal name) el=(Pascal name)}}
{{/with}}
`;
