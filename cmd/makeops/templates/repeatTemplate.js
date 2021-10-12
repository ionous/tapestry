// repeatsPartial.js
// note: the non-repeating callbacks used by these functions are expected to be defined elsewhere
// ex. in flow, or prim, etc.
'use strict';
module.exports =`
type {{name}}{{mod}}_Slice []{{el}}

func (op* {{name}}{{mod}}_Slice) GetSize() int    { return len(*op) }
func (op* {{name}}{{mod}}_Slice) SetSize(cnt int) { (*op) = make({{name}}{{mod}}_Slice, cnt) }

func {{name}}{{mod}}_Repeats_Marshal(n jsn.Marshaler, vals *[]{{el}}) {
  if n.RepeatValues({{name}}_Type, (*{{name}}{{mod}}_Slice)(vals)) {
    for i := range *vals {
      {{name}}{{mod}}_Marshal(n, &(*vals)[i])
    }
    n.EndValues()
  }
  return
}
`;
