// strTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  Str string
}

func (op *{{Pascal name}}) String()(ret string) {
  return op.Str
}

{{#each (Choices @this)}}
const {{Pascal ../name}}_{{Pascal this.token}}= "{{this.token}}";
{{/each}}


{{>spec spec=this}}
`;
