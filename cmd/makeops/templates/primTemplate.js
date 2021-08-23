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
{{/each~}}

{{>spec}}

func {{Pascal name}}_Exists(val *{{Pascal name}}) bool {
  var zero {{Pascal name}}
  return val.{{Pascal uses}} != zero.{{Pascal uses}}
}
{{#if (OverrideOf name)}}
func {{Pascal name}}_Override_Exists(val *{{OverrideOf name}}) bool {
  var zero {{OverrideOf name}}
  return *val != zero
}
{{/if}}
{{/with}}
{{#if marshal}}
{{>sig type}}
{{~#if (OverrideOf type.name)}}
{{>override type fmt="Compact"}}
{{>override type fmt="Detailed"}}
{{>repeat name=(Pascal type.name) mod="_Override" el=(OverrideOf type.name)}}
{{/if}}
{{>primCompact type}}
{{>primDetails type}}
{{/if}}

{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
