// specPartial.js
'use strict';
module.exports =`func (*{{Pascal name}}) Compose() composer.Spec {
  return composer.Spec{
    Name: "{{name}}",
{{#if (LedeName this)}}
    Lede: "{{LedeName this}}",
{{/if}}
{{#if (IsStr name)}}
{{#unless (IsClosed this)}}
    OpenStrings: true,
{{/unless}}
{{#if (Choices this)}}
    Strings: []string{
      {{#each (Choices @this)~}}"{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
    },
{{/if}}
{{/if}}
  }
}
`;
