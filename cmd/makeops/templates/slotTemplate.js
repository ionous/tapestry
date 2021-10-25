// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const {{Pascal name}}_Type = "{{name}}";
{{#if ../marshal}}
var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal

type {{Pascal name}}_Slot struct { ptr *{{Pascal name}} }

func (At {{Pascal name}}_Slot) GetType() string { return {{Pascal name}}_Type }
func (at {{Pascal name}}_Slot) GetSlot() (interface{}, bool) { return at.ptr, at.ptr != nil }
func (at {{Pascal name}}_Slot) SetSlot(v interface{}) (okay bool) {
  (*at.ptr), okay = v.({{Pascal name}})
  return
}

func {{Pascal name}}_Marshal(m jsn.Marshaler, ptr *{{Pascal name}}) (err error) {
{{#if (IsCustom name)}}
  if fn, exists := m.CustomizedMarshal({{Pascal name}}_Type); exists {
    err = fn(m, ptr)
  } else {
    err = {{Pascal name}}_DefaultMarshal(m, ptr)
  }
  return
}
func {{Pascal name}}_DefaultMarshal(m jsn.Marshaler, ptr *{{Pascal name}}) (err error) {
{{/if}}
{{~#if (IsPositioned this)}}
  m.SetCursor(ptr.At.Offset)
{{/if}}
  if err = m.MarshalBlock({{Pascal name}}_Slot{ptr}); err == nil {
    if e := (*ptr).(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
      m.Error(e)
    }
    m.EndBlock()
  }
  return
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
