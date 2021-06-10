// swapTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} swaps between various options
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  Opt interface{}
}

{{>spec spec=this}}

{{#each with.params}}
const {{Pascal ../name}}_{{Pascal @key}}= "{{@key}}";
{{/each}}

`;
