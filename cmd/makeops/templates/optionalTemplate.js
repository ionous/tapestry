// optionalTemplate.js
'use strict';
module.exports =`{{#if optional}}{{#with type}}
func {{Pascal name}}_Detailed_Optional_Marshal(n jsonexp.Context, val **{{Pascal name}}) (ret []byte, err error) {
  if ptr:= *val; ptr != nil {
    ret, err = {{Pascal name}}_Detailed_Marshal(n, ptr)
  }
  return
}
func {{Pascal name}}_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **{{Pascal name}}) (err error) {
  if len(b) > 0 {
    var el {{Pascal name}}
    if e := {{Pascal name}}_Detailed_Unmarshal(n, b, &el); e!= nil {
      err = e
    } else {
      *out = &el
    }
  }
  return
}
{{/with}}{{/if}}`;
