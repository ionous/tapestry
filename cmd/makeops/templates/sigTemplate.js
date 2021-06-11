// sigTemplate.js
module.exports = `
func (op* {{Pascal name}}) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
  return {{Pascal name}}_Detailed_Marshal(n, op)
}
func (op *{{Pascal name}}) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
  return {{Pascal name}}_Detailed_Unmarshal(n, b, op)
}
`
