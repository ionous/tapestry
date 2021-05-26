// strTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} string


func (*{{Pascal name}}) Choices() (choices map[string]string) {
  return map[string]string{
    {{#each (Choices @this)~}}"{{this.token}}": "{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
  }
}

{{>spec spec=this}}
`;
