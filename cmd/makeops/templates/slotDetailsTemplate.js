// slotTemplate.js
'use strict';
module.exports = `
var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal

func {{Pascal name}}_Marshal(n jsn.Marshaler, ptr *{{Pascal name}}) {
  if slat := *ptr; slat != nil {
    slat.(jsn.Marshalee).Marshal(n)
  }
  return
}
`;
