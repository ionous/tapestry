// runTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} {{desc.short}}
type {{Pascal name}} struct {
{{#each (ParamsOf this)}}
  {{Pascal @key}} {{TypeOf this}} \`if:"
  {{~#if (IsInternal label)}}internal{{else}}label={{LabelOf @key this @index}}{{/if}}
  {{~#if optional}},optional{{/if}}\
  {{~#if (OverrideOf type)}},type={{type}}{{/if}}"\`
{{/each}}
}

{{>spec~}}
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
const {{Pascal ../name}}_{{Pascal @key}}= "{{@key}}";
{{~/unless}}{{~/each}}
{{/with}}
{{#if marshal}}
{{>sig type}}
{{>flowDetails type}}
{{/if}}
{{#if repeats}}{{>repeat type}}{{/if}}
{{#if  optional}}{{>option type}}{{/if}}
`;
