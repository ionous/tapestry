// slotTemplate.js
'use strict';
module.exports = `
var {{Pascal name}}_Compact_Optional_Marshal = {{Pascal name}}_Compact_Marshal
var {{Pascal name}}_Compact_Optional_Unmarshal = {{Pascal name}}_Compact_Unmarshal

func {{Pascal name}}_Compact_Marshal(n jsonexp.Context, ptr *{{Pascal name}}) (ret []byte ,err error) {
  if slat := *ptr; slat != nil {
    ret, err = slat.(jsonexp.CompactMarshaler).MarshalCompact(n)
  }
  return
}

func {{Pascal name}}_Compact_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  if ptr, e := jsonexp.UnmarshalCompactSlot(n, b); e != nil {
    err =e
  } else if store, ok := ptr.({{Pascal name}}); !ok && ptr != nil {
    err = errutil.Fmt("couldnt store %T into %s", ptr, {{Pascal name}}_Type)
  } else {
    (*out) = store
  }
  return
}
`;
