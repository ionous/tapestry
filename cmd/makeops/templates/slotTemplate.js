// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const {{Pascal name}}_Type = "{{name}}";
{{#if ../marshal}}
var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal

type {{Pascal name}}_Slot struct { ptr *{{Pascal name}} }

func (At {{Pascal name}}_Slot) GetType() string { return {{Pascal name}}_Type }
func (at {{Pascal name}}_Slot) HasSlot() bool { return at.ptr != nil }
func (at {{Pascal name}}_Slot) SetSlot(v interface{}) (okay bool) {
  (*at.ptr), okay = v.({{Pascal name}})
  return
}

func {{Pascal name}}_Marshal(n jsn.Marshaler, ptr *{{Pascal name}}) (okay bool) {
{{#if (IsCustom name)}}
  if fn, exists := n.CustomizedMarshal({{Pascal name}}_Type); exists {
    okay = fn(n, ptr)
  } else {
    okay = {{Pascal name}}_DefaultMarshal(n, ptr)
  }
  return
}
func {{Pascal name}}_DefaultMarshal(n jsn.Marshaler, ptr *{{Pascal name}}) (okay bool) {
{{/if}}
{{~#if (IsPositioned this)}}
  n.SetCursor(ptr.At.Offset)
{{/if}}
  if okay = n.MarshalBlock({{Pascal name}}_Slot{ptr}); okay {
    (*ptr).(jsn.Marshalee).Marshal(n)
    n.EndBlock()
  }
  return
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
