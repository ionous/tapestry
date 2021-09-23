// swapDetailsTemplate.js
'use strict';
module.exports = `
func {{Pascal name}}_Marshal(n jsn.Marshaler, val *{{Pascal name}}){
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
`;
