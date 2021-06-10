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

{{>spec spec=this}}

{{#if ../marshal}}
func (op *{{Pascal name}}) MarshalJSON() ([]byte, error) {
  return json.Marshal(jsonexp.String{
{{#if (IsPositioned this)}}
    Id: op.At.Offset,
{{/if}}
    Type:  Type_{{Pascal name}},
    Value: op.Str,
  })
}

func (op *{{Pascal name}}) UnmarshalJSON(b []byte) (err error) {
  var d jsonexp.String;
  if e := json.Unmarshal(b, &d); e != nil {
    err= e
  } else {
{{#if (IsPositioned this)}}
    op.At.Offset= d.Id;
{{/if}}
    op.Str= d.Value;
  }
  return
}
{{/if}}

{{#each (Choices @this)}}
const {{Pascal ../name}}_{{Pascal this.token}}= "{{this.token}}";
{{/each~}}

{{/with}}
`;
