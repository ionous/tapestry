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
  if *val != zero {
    {{Pascal name}}_Unboxed_Marshal(n, val)
  }
}

func {{Pascal name}}_Unboxed_Marshal(n jsn.Marshaler, val *{{Unboxed name}}) {
  {{Pascal name}}_Marshal(n, &{{Pascal name}}{jsn.Box{{Pascal (Unboxed name)}}(val)})
}
{{/if}}

func {{Pascal name}}_Optional_Marshal(n jsn.Marshaler, val *{{Pascal name}}) {
  var zero {{Pascal name}}
  if val.{{Pascal uses}} != zero.{{Pascal uses}} {
     {{Pascal name}}_Marshal(n, val)
  }
}

func {{Pascal name}}_Marshal{{Custom name}}(n jsn.Marshaler, val *{{Pascal name}}) {
{{~#if (IsPositioned this)}}
  n.SetCursor(val.At.Offset)
{{/if}}
{{#if (IsEnumerated this)}}
  n.WriteChoice({{Pascal name}}_Type, val)
{{else}}
  n.WriteValue({{Pascal name}}_Type, val.{{Pascal uses}})
{{/if}}
}

{{/if}}
{{>repeat name=(Pascal name) el=(Pascal name)}}
{{/with}}
`;
