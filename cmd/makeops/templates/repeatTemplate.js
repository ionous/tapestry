// repeatsPartial.js
// note: the non-repeating callbacks used by these functions are expected to be defined elsewhere
// ex. in flow, or prim, etc.
'use strict';
module.exports =`
type {{name}}{{mod}}_Slice []{{el}}

func (op* {{name}}{{mod}}_Slice) GetType() string { return {{name}}_Type }
func (op* {{name}}{{mod}}_Slice) GetSize() int    { return len(*op) }
func (op* {{name}}{{mod}}_Slice) SetSize(cnt int) { (*op) = make({{name}}{{mod}}_Slice, cnt) }

func {{name}}{{mod}}_Repeats_Marshal(m jsn.Marshaler, vals *[]{{el}}) (err error) {
  if err = m.MarshalBlock((*{{name}}{{mod}}_Slice)(vals)); err == nil {
    for i := range *vals {
      if e := {{name}}{{mod}}_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
        m.Error(errutil.New(e, "in slice at", i))
      }
    }
    m.EndBlock()
  }
  return
}
`;
