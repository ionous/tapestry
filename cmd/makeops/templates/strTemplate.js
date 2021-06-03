// strTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
  At    reader.Position \`if:"internal"\`
  Str string
}

func (op *{{Pascal name}}) String()(ret string) {
  if s := op.Str; s != "$EMPTY" {
    ret = s
  }
  return
}

{{#each (Choices @this)}}
const {{Pascal ../name}}_{{Pascal this.token}}= "{{this.token}}";
{{/each}}

func (*{{Pascal name}}) Choices() (choices map[string]string) {
  return map[string]string{
    {{#each (Choices @this)~}}
    {{Pascal ../name}}_{{Pascal this.token}}: "{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
  }
}

{{>spec spec=this}}
`;
