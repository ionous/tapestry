// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const {{Pascal name}}_Type = "{{name}}";
{{#if ../marshal}}
var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal

type {{Pascal name}}_Slot struct { Value *{{Pascal name}} }

func (at *{{Pascal name}}_Slot) Marshal(m jsn.Marshaler) (err error) {
  if err = m.MarshalBlock(at); err == nil {
    if a, ok := at.GetSlot(); ok {
      if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
        m.Error(e)
      }
    }
    m.EndBlock()
  }
  return
}
func (at *{{Pascal name}}_Slot) GetType() string { return {{Pascal name}}_Type }
func (at *{{Pascal name}}_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at *{{Pascal name}}_Slot) SetSlot(v interface{}) (okay bool) {
  (*at.Value), okay = v.({{Pascal name}})
  return
}

func {{Pascal name}}_Marshal(m jsn.Marshaler, ptr *{{Pascal name}}) (err error) {
{{~#if (IsPositioned this)}}
  m.SetCursor(ptr.At.Offset)
{{/if}}
  slot := {{Pascal name}}_Slot{ptr}
  return slot.Marshal(m)
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
