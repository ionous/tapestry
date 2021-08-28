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
const {{Pascal ../name}}_Field_{{Pascal @key}} = "{{@key}}";
{{~/unless}}{{~/each}}
{{#if ../marshal}}
{{>sig}}
{{#unless (NoHelpers name)}}
{{>repeat name=(Pascal name) el=(Pascal name)}}
{{/unless}}
{{>optional mod="Compact"}}
{{> (compactFlow this)}}
{{>optional mod="Detailed"}}
{{>flowDetails}}
{{/if}}
{{/with}}
`;
