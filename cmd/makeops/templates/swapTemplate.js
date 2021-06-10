// swapTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} swaps between various options
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  Opt interface{}
}

{{#each with.params}}
const {{Pascal ../name}}_{{Pascal @key}}= "{{@key}}";
{{/each}}

{{>spec spec=this}}

func (op* {{Pascal name}}) MarshalJSON() (ret []byte, err error) {
  if val, e:= op.MarshalChoice(); e!= nil {
    err=e
  } else {
    ret, err= json.Marshal(map[string]interface{}{
{{~#if (IsPositioned this)}}
      "id": op.At.Offset,{{/if}}
      "type": "{{name}}",
      "value": val,
    })
  }
  return
}

func (op* {{Pascal name}}) MarshalChoice() (ret map[string]json.RawMessage, err error) {
  if kid, e:= json.Marshal(op.Opt); e!= nil {
    err= e
  } else if pick, ok:= op.GetChoice(); !ok {
    err= errutil.Fmt("unknown choice %T in %T", op.Opt, op)
  } else if (len(pick) > 0) {
    ret= map[string]json.RawMessage{
      pick: kid,
    }
  }
  return
}

func (op* {{Pascal name}}) GetChoice() (ret string, okay bool) {
  switch op.Opt.(type) {
  case nil:
    okay= true
{{#each with.params}}
  case *{{TypeOf this}}:
    ret, okay= {{Pascal ../name}}_{{Pascal @key}}, true
{{/each}}
  }
  return
}
{{/with}}
`;
