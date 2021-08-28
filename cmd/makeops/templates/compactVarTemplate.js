// compactVarTemplate.js
'use strict';
module.exports = `
func GetVar_Compact_Marshal(n jsonexp.Context, val *GetVar) (ret []byte, err error) {
  // custom serialization collapses variable commands into strings starting with @
  if str := val.Name.Str; len(str) > 0 && str[0] == '@' {
    err = errutil.New("serialization doesn't support variables names starting with @")
  } else {
    ret, eerr = json.Marshal("@" + str)
  }
  return
}
func GetVar_Compact_Unmarshal(n jsonexp.Context, b []byte, out *GetVar) (err error) {
  return
}
`;
