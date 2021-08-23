// slotTemplate.js
'use strict';
module.exports = `
var {{Pascal name}}_Detailed_Optional_Marshal = {{Pascal name}}_Detailed_Marshal
var {{Pascal name}}_Detailed_Optional_Unmarshal = {{Pascal name}}_Detailed_Unmarshal

func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, ptr *{{Pascal name}}) (ret []byte,err error) {
  if slat := *ptr; slat != nil {
    ret, err = slat.(jsonexp.DetailedMarshaler).MarshalDetailed(n)
  }
  return
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  if ptr, e := jsonexp.UnmarshalDetailedSlot(n, b); e != nil {
    err =e
  } else if store, ok := ptr.({{Pascal name}}); !ok && ptr != nil {
    err = errutil.Fmt("couldnt store %T into %s", ptr, {{Pascal name}}_Type)
  } else {
    (*out) = store
  }
  return
}
`;
