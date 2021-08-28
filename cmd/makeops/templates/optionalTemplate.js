// runTemplate.js
'use strict';
module.exports = `
func {{Pascal name}}_{{mod}}_Optional_Marshal(n jsonexp.Context, val **{{Pascal name}}) (ret []byte,err error) {
  if *val != nil {
    ret, err= {{Pascal name}}_{{mod}}_Marshal(n, *val)
  }
  return
}
func {{Pascal name}}_{{mod}}_Optional_Unmarshal(n jsonexp.Context, b []byte, out **{{Pascal name}}) (err error) {
  if len(b) > 0 {
    var val {{Pascal name}}
    if e:= {{Pascal name}}_{{mod}}_Unmarshal(n, b, &val); e!= nil {
      err = e
    } else {
      *out = &val
    }
  }
  return
}`
