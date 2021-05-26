// runTemplate.js
'use strict';
module.exports =`
// {{Pascal name}} {{desc.short}}
type {{Pascal name}} struct {
{{#each with.params}}
  {{NameOf @key this}} {{TypeOf this}} \`if:"
  {{~#if (IsInternal label)}}internal{{else}}label={{label}}{{/if}}
  {{~#if optional}},optional{{/if}}"\`
{{/each}}
}

{{>spec spec=this}}
`;
