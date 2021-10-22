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

func (op* {{Pascal name}}) GetType() string { return {{Pascal name}}_Type }

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

func (op* {{Pascal name}}) SetChoice(c string) (ret interface{}, okay bool) {
  switch c {
  case "":
    op.Opt, okay = nil, true
{{#each with.params}}
  case {{Pascal ../name}}_{{Pascal @key}}_Opt:
    opt := new({{TypeOf this}})
    op.Opt, ret, okay = opt, opt, true
{{/each}}
  }
  return
}

{{#if ../marshal}}
{{>sig}}
func {{Pascal name}}_Marshal(n jsn.Marshaler, val *{{Pascal name}}){
{{#if (IsCustom name)}}
  if fn, ok := n.CustomizedMarshal({{Pascal name}}_Type); ok {
    fn(n, val)
  } else {
    {{Pascal name}}_DefaultMarshal(n, val)
  }
  return
}
func {{Pascal name}}_DefaultMarshal(n jsn.Marshaler, val *{{Pascal name}}) {
{{/if}}
{{~#if (IsPositioned this)}}
  n.SetCursor(val.At.Offset)
{{/if}}
  if n.MarshalBlock(val) {
    val.Opt.(jsn.Marshalee).Marshal(n)
    n.EndBlock()
  }
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
