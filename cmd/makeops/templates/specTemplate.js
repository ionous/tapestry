// specPartial.js
'use strict';
module.exports =`func (*{{Pascal name}}) Compose() composer.Spec {
  return composer.Spec{
    Name: {{Pascal name}}_Type,
    Uses: composer.Type_{{Pascal uses}},
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
      {{#each (Choices @this)~}}{{Pascal ../name}}_{{Pascal this.token}}_Opt,{{#unless @last}} {{/unless}}{{/each}}
     },
    Swaps: []interface{} {
{{#each with.params}}
      (*{{TypeOf this}})(nil),
{{/each}}
    },
{{/if}}
  }
}

{{#if (Uses name "str")}}{{#if (Choices this)}}
func (op *{{Pascal name}}) SetEnum(kv string) {
  composer.SetEnum(op, kv, &op.Str)
}
func (op *{{Pascal name}}) GetEnum() (retKey string, retVal string) {
  return composer.GetEnum(op, op.Str)
}
{{/if}}{{/if}}
const {{Pascal name}}_Type = "{{name}}"
`;
