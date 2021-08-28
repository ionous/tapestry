// compactTextTemplate.js
'use strict';
module.exports = `
func TextValue_Compact_Marshal(n jsonexp.Context, val *TextValue) ([]byte, error) {
  // custom serialization to avoid conflicts with @variables
  str := val.Text
  if len(str) > 0 && str[0] == '@' {
    str = "@"+ str
  }
  return json.Marshal(str)
}
func {{Pascal name}}_Compact_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  return
}
`;
