// numTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} requires a user-specified number.
type {{Pascal name}} struct {
  Value float64
}

{{#if ../marshal}}
func (op *{{Pascal name}}) MarshalJSON() ([]byte, error) {
  return json.Marshal(map[string]interface{}{
{{~#if (IsPositioned this)}}
    "id": op.At.Offset,{{/if}}
    "type": "{{name}}",
    "value": op.Value,
  })
}
{{/if}}


{{>spec spec=this}}
{{/with}}
`;
