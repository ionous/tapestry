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
const {{Pascal ../name}}_{{Pascal @key}}_Opt= "{{@key}}";
{{/each}}

{{>spec}}

func (op* {{Pascal name}}) GetChoice() (ret string, okay bool) {
  switch op.Opt.(type) {
  case nil:
    okay = true
{{#each with.params}}
  case *{{TypeOf this}}:
    ret, okay = {{Pascal ../name}}_{{Pascal @key}}_Opt, true
{{/each}}
  }
  return
}

{{~#if ../marshal}}
{{>sig}}
func {{Pascal name}}_Marshal{{Custom name}}(n jsn.Marshaler, val *{{Pascal name}}){
  if pick, ok := val.GetChoice(); ok {
    if slat := val.Opt; len(pick) > 0 {
{{~#if (IsPositioned this)}}
      n.SetCursor(val.At.Offset)
{{/if}}
      n.PickValues({{Pascal name}}_Type, pick)
      /* */ slat.(jsn.Marshalee).Marshal(n)
      n.EndValues()
    }
  }
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
