// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const {{Pascal name}}_Type = "{{name}}";
{{#if ../marshal}}
var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal

type {{Pascal name}}_Slot struct { ptr *{{Pascal name}} }

func (at {{Pascal name}}_Slot) HasSlot() bool { return at.ptr != nil }
func (at {{Pascal name}}_Slot) SetSlot(v interface{}) (okay bool) {
  (*at.ptr), okay = v.({{Pascal name}})
  return
}

func {{Pascal name}}_Marshal(n jsn.Marshaler, ptr *{{Pascal name}}) {
{{#if (IsCustom name)}}
  if fn, ok := n.CustomizedMarshal({{Pascal name}}_Type); ok {
    fn(n, ptr)
  } else {
    {{Pascal name}}_DefaultMarshal(n, ptr)
  }
  return
}
func {{Pascal name}}_DefaultMarshal(n jsn.Marshaler, ptr *{{Pascal name}}) {
{{/if}}
{{~#if (IsPositioned this)}}
  n.SetCursor(ptr.At.Offset)
{{/if}}
  if ok := n.SlotValues({{Pascal name}}_Type, {{Pascal name}}_Slot{ptr}); ok {
    (*ptr).(jsn.Marshalee).Marshal(n)
    n.EndValues()
  }
  return
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
