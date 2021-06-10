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
  return json.Marshal(jsonexp.Float{
{{#if (IsPositioned this)}}
    Id: op.At.Offset,
{{/if}}
    Type:  Type_{{Pascal name}},
    Value: op.Value,
  })
}

func (op *{{Pascal name}}) UnmarshalJSON(b []byte) (err error) {
  var d jsonexp.Float;
  if e := json.Unmarshal(b, &d); e != nil {
    err= e
  } else {
{{#if (IsPositioned this)}}
    op.At.Offset= d.Id;
{{/if}}
    op.Value= d.Value;
  }
  return
}
{{/if}}


{{>spec spec=this}}
{{/with}}
`;
