// repeatsPartial.js
// note: the non-repeating callbacks used by these functions are expected to be defined elsewhere
// ex. in flow, or prim, etc.
'use strict';
module.exports =`
func {{name}}{{mod}}_Repeats_Marshal(n jsn.Marshaler, vals *[]{{el}}) {
  if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
    n.RepeatValues(cnt)
    for _, el := range *vals {
      {{name}}{{mod}}_Marshal(n, &el)
    }
    n.EndValues()
  }
  return
}
`;
