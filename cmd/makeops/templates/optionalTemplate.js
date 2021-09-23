// runTemplate.js
'use strict';
module.exports = `
func {{Pascal name}}_Optional_Marshal(n jsn.Marshaler, val **{{Pascal name}}) {
  if *val != nil {
    {{Pascal name}}_Marshal(n, *val)
  }
}
`
