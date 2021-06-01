// strTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
  Str string
}

func (op *{{Pascal name}}) String()(ret string) {
  if s := op.Str; s != "$EMPTY" {
    ret = s
  }
  return
}

func (*{{Pascal name}}) Choices() (choices map[string]string) {
  return map[string]string{
    {{#each (Choices @this)~}}"{{this.token}}": "{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
  }
}

{{>spec spec=this}}
`;
