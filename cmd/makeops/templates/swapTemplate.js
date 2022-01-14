// swapTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} swaps between various options
type {{Pascal name}} struct {
  Choice string
  Value  interface{}
}

var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal

{{#each with.params}}
const {{Pascal ../name}}_{{Pascal @key}}_Opt= "{{@key}}";
{{/each}}

{{>spec}}

func (op* {{Pascal name}}) GetType() string { return {{Pascal name}}_Type }

func (op* {{Pascal name}}) GetSwap() (string, interface{}) {
  return op.Choice, op.Value
}

func (op* {{Pascal name}}) SetSwap(c string) (okay bool) {
  switch c {
  case "":
    op.Choice, op.Value = c, nil
    okay = true
{{#each with.params}}
  case {{Pascal ../name}}_{{Pascal @key}}_Opt:
    op.Choice, op.Value = c, new({{TypeOf this}})
    okay = true
{{/each}}
  }
  return
}

{{#if ../marshal}}
{{>sig}}
func {{Pascal name}}_Marshal(m jsn.Marshaler, val *{{Pascal name}}) (err error) {
  if err = m.MarshalBlock(val); err == nil {
    if _, ptr := val.GetSwap(); ptr != nil {
      if e := ptr.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
        m.Error(e)
      }
    }
    m.EndBlock()
  }
  return
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
