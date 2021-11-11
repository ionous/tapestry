// repeatsPartial.js
// note: the non-repeating callbacks used by these functions are expected to be defined elsewhere
// ex. in flow, or prim, etc.
'use strict';
module.exports =`
type {{name}}{{mod}}_Slice []{{el}}

func (op* {{name}}{{mod}}_Slice) GetType() string {  return {{name}}_Type }

func (op *{{name}}{{mod}}_Slice) Marshal(m jsn.Marshaler) error {
  return {{name}}{{mod}}_Repeats_Marshal(m, (*[]{{el}})(op))
}

func (op* {{name}}{{mod}}_Slice) GetSize() (ret int) {
  if els:= *op; els != nil {
    ret = len(els)
  } else {
    ret = -1
   }
   return
}

func (op* {{name}}{{mod}}_Slice) SetSize(cnt int) {
  var els []{{el}}
  if cnt >= 0 {
    els = make({{name}}{{mod}}_Slice, cnt)
  }
  (*op) = els
}

func (op* {{name}}{{mod}}_Slice) MarshalEl(m jsn.Marshaler, i int) error {
  return {{name}}{{mod}}_Marshal(m, &(*op)[i])
}

func {{name}}{{mod}}_Repeats_Marshal(m jsn.Marshaler, vals *[]{{el}}) error {
  return jsn.RepeatBlock(m, (*{{name}}{{mod}}_Slice)(vals))
}

func {{name}}{{mod}}_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]{{el}}) (err error) {
  if *pv != nil || !m.IsEncoding() {
    err = {{name}}{{mod}}_Repeats_Marshal(m, pv)
  }
  return
}
`;
