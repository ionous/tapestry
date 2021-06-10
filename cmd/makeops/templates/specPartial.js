// specPartial.js
'use strict';
module.exports =`func (*{{Pascal name}}) Compose() composer.Spec {
  return composer.Spec{
    Name: "{{name}}",
    Uses: "{{uses}}",
{{#if (LedeName this)}}
    Lede: "{{LedeName this}}",
{{/if}}
{{#if (Uses name "str")}}
{{#unless (IsClosed this)}}
    OpenStrings: true,
{{/unless}}
{{#if (Choices this)}}
    Choices: []string {
      {{#each (Choices @this)~}}{{Pascal ../name}}_{{Pascal this.token}},{{#unless @last}} {{/unless}}{{/each}}
     },
    Strings: []string{
      {{#each (Choices @this)~}}"{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
    },
{{/if}}
{{/if}}
{{#if (Uses name "swap")}}
    Choices: []string {
      {{#each (Swaps @this)~}}{{Pascal ../name}}_{{Pascal this.token}},{{#unless @last}} {{/unless}}{{/each}}
     },
    Swaps: []interface{} {
{{#each with.params}}
      (*{{TypeOf this}})(nil),
{{/each}}
    },
{{/if}}
  }
}
`;
