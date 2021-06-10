// strTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  Str string
}

func (op *{{Pascal name}}) String()(ret string) {
  return op.Str
}

{{#if ../marshal}}
func (op *{{Pascal name}}) MarshalJSON() ([]byte, error) {
  return json.Marshal(map[string]interface{}{
{{~#if (IsPositioned this)}}
    "id": op.At.Offset,{{/if}}
    "type": "{{name}}",
    "value": op.Str,
  })
}
{{/if}}

{{#each (Choices @this)}}
const {{Pascal ../name}}_{{Pascal this.token}}= "{{this.token}}";
{{/each~}}

{{>spec spec=this}}
{{/with}}
`;
