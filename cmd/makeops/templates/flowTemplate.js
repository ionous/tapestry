// runTemplate.js
'use strict';
module.exports =`
// {{Pascal name}} {{desc.short}}
type {{Pascal name}} struct {
{{~#each (ParamsOf this)}}
  {{NameOf @key this}} {{TypeOf this}} \`if:"
  {{~#if (IsInternal label)}}internal{{else}}label={{LabelOf label}}{{/if}}
  {{~#if optional}},optional{{/if}}"\`
{{/each}}
}

{{>spec spec=this}}
`;
