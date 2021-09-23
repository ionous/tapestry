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
{{/with}}
{{#if marshal}}
{{>sig type}}
{{~#if (OverrideOf type.name)}}
{{>override type}}
{{>repeat name=(Pascal type.name) mod="_Override" el=(OverrideOf type.name)}}
{{/if}}
{{>primDetails type}}
{{/if}}

{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
