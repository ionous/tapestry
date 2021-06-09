// runTemplate.js
'use strict';
module.exports =`
// {{Pascal name}} {{desc.short}}
type {{Pascal name}} struct {
{{~#each (ParamsOf this)}}
  {{ParamNameOf @key this}} {{TypeOf this}} \`if:"
  {{~#if (IsInternal label)}}internal{{else}}label={{LabelOf @key this @index}}{{/if}}
  {{~#if optional}},optional{{/if}}\
  {{~#if (Override this)}},type={{Override this}}{{/if}}"\`
{{/each}}
}
{{#each this.with.slots}}
var _ {{ScopedNameOf this}} = (*{{Pascal ../name}})(nil)
{{/each}}
{{>spec spec=this}}
`;
